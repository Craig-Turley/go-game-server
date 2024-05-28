package models

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Canvas struct {
    height float32
    width  float32
}

type GameState  string
type PlayerTurn string

const (
    INLOBBY  = "INLOBBY"
    PLAYING  = "PLAYING"
    FINISHED = "FINISHED"

    PLAYER1 = "PLAYER1"
    PLAYER2 = "PLAYER2"
)

type Queue interface {
    Add(cl *Client)
}

type PrivateQueue struct {
    Manager   *PrivateGameManager
    Clients   []*Client
    Id        string
    Host      *Client
    State     GameState
    Turn      PlayerTurn
    ball      Ball
    score     Score
}

func (p *PrivateQueue) Add(cl *Client) {
    p.Clients = append(p.Clients, cl)
}

func (p *PrivateQueue) MonitorStart() {
    for {
        _, m, err := p.Host.Ws.ReadMessage()
        if err != nil {
            log.Println(err)
            p.cleanUp()
            return
        }

        log.Println(string(m))
        if string(m) == "GameStart" {
            for i, c := range p.Clients {
                gd := GameStart {
                    Message: true,
                    Player: i + 1,
                }
                c.SetPosition(i + 1)
                if err := c.Ws.WriteJSON(gd); err != nil {
                    log.Println(err)
                    p.cleanUp()
                    return
                }
            }
            go p.startGame()
            log.Println("MonitorStart() returning...")
            return
        }
    }
}

func(p *PrivateQueue) startGame() {
    gameEnd := make(chan bool)
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    for i := 3; i > 0; i-- {
        for _, c := range p.Clients {
            if err := c.Ws.WriteMessage(1, []byte(fmt.Sprintf("%d", i))); err != nil {
                log.Println(err)
                p.cleanUp()
                cancel()
                return
            }
        }
        time.Sleep(1 * time.Second)
    }
    for _, c := range p.Clients {
        if err := c.Ws.WriteMessage(1, []byte("Start")); err != nil {
            log.Println(err)
            cancel()
            p.cleanUp()
            return
        }
        go p.readAndUpdateClientGameState(ctx, c, gameEnd)
    }
    go p.updateBall(ctx, gameEnd)

    // monitors game end
    // Go routine hangs until signal is sent via children that game has finsihed
    <-gameEnd
    cancel()
    log.Println("Game end ID: " + p.Id)
    payload := GameDataSend {
        PaddleY:        0,
        BallY:          p.ball.Y,
        BallX:          p.ball.X,
        PlayerOneScore: p.score.PlayerOne,
        PlayerTwoScore: p.score.PlayerTwo,
    }

    for _, c := range p.Clients {
        c.Ws.WriteJSON(payload)
    }

    // signal game end to every client
    // clean up process
    p.cleanUp()
}

func (p *PrivateQueue) readAndUpdateClientGameState(ctx context.Context, sender *Client, gameEnd chan bool) {
    for {
        select {
        case <-ctx.Done():
            log.Println("readAndUpdateClientGameState() returning...")
            return
        default:
            childCtx, cancel := context.WithTimeout(ctx, 1 * time.Second)
            defer cancel()
            ch := make(chan GameDataReceive, 1)
            go func(ctx context.Context, ch chan GameDataReceive) {
                var data GameDataReceive
                err := sender.Ws.ReadJSON(&data);
                if err != nil {
                    handleError(err, gameEnd)
                }
                ch <- data
            }(childCtx, ch)

            select {
            case data := <-ch:
                sender.Paddle.Set(data.PaddleX, data.PaddleY)

                payload := GameDataSend {
                    PaddleY:        data.PaddleY,
                    BallY:          p.ball.Y,
                    BallX:          p.ball.X,
                    PlayerOneScore: p.score.PlayerOne,
                    PlayerTwoScore: p.score.PlayerTwo,
                }

                for _, c := range p.Clients {
                    if c == sender { continue }

                    c.Ws.WriteJSON(payload)
                }
            case <- childCtx.Done():
                log.Println("readAndUpdateClientGameState() returning...")
                return
            }

        }
    }
}

var canvas = Canvas{width: 1400, height: 1000}

func (p *PrivateQueue) updateBall(ctx context.Context, gameEnd chan bool) {
    interval := time.Second / 60
    ball := &p.ball
    p.ball.Randomize(canvas.height)
    for {
        select {
        case <-ctx.Done():
            log.Println("updateBall() returning...")
            return
        default:
            start := time.Now()

            // change x depending on speed and direction
            switch ball.DirectionX {
            case LEFT:
                ball.X -= ball.Speed
                break
            case RIGHT:
                ball.X += ball.Speed
                break
            }

            // change y depending on speed and direction
            switch ball.DirectionY {
            case UP:
                ball.Y -= ball.Speed / 1.5
                break
            case DOWN:
                ball.Y += ball.Speed / 1.5
                break
            }

            // handle ball out of bounds collisions

            switch {
            case ball.X <= 0:
                p._resetTurn(2, gameEnd)
            case ball.X >= canvas.width - ball.Width():
                p._resetTurn(1, gameEnd)
            case ball.Y <= 0:
                ball.DirectionY = DOWN
            case ball.Y >= canvas.height - ball.Height():
                ball.DirectionY = UP
            }

            for _, c := range p.Clients {
                if p.ball.X - p.ball.Width() <= c.Paddle.X && ball.X >= c.Paddle.X - c.Paddle.Width() {
                    if p.ball.Y <= c.Paddle.Y + c.Paddle.Height() && p.ball.Y + p.ball.Height() >= c.Paddle.Y {
                        switch c.PlayerPosition {
                        case 1:
                            p.ball.X = c.Paddle.X + p.ball.Width()
                            p.ball.DirectionX = RIGHT
                        case 2:
                            p.ball.X = c.Paddle.X - p.ball.Width()
                            p.ball.DirectionX = LEFT
                        }
                    }
                }
            }
            time.Sleep(interval - time.Since(start))
        }
    }
}

func (p *PrivateQueue) _resetTurn(victor int, gameEnd chan bool) {

    switch victor {
    case 1:
        p.score.PlayerOne += 1
        if p.score.PlayerOne == 10 {
            gameEnd<-true
        }
        break
    case 2:
        p.score.PlayerTwo += 1
        if p.score.PlayerTwo == 10 {
            gameEnd<-true
        }
        break
    }

    p.ball.Reset()
    p.ball.Randomize(canvas.height)

    time.Sleep(1 * time.Second)
}

func handleError(err error, gameEnd chan bool) {
    if err != nil {
        if ce, ok := err.(*websocket.CloseError); ok {
            switch ce.Code {
            case websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived:
                gameEnd <- true
            }
            log.Println(err)
        }
    }
}



func (p *PrivateQueue) cleanUp() {
    for _, c := range p.Clients {
        c.Ws.WriteMessage(1, []byte("Game has ended"))
        c.Ws.Close()
        p.Manager.CleanUp(p.Id)
    }
}
