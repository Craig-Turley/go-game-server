package models

import (
	// "fmt"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Canvas struct {
    height int
    width  int
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
    Clients []*Client
    Id      string
    Host    *Client
    State   GameState
    Turn    PlayerTurn
}

func (p *PrivateQueue) Add(cl *Client) {
    p.Clients = append(p.Clients, cl)
}

func (p *PrivateQueue) MonitorStart() {
    for {
        _, m, err := p.Host.Ws.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }

        log.Println(string(m))
        if string(m) == "GameStart" {
            for i, c := range p.Clients {
                gd := GameStart {
                    Message: true,
                    Player: i + 1,
                }
                log.Println(gd)
                if err := c.Ws.WriteJSON(gd); err != nil {
                    log.Println(err)
                    return
                }
            }
            go p.startGame()
            return
        }
    }
}

func(p *PrivateQueue) startGame() {
    gameEnd := make(chan bool)
    go p.monitorEnd(gameEnd)
    for i := 3; i > 0; i-- {
        for _, c := range p.Clients {
            if err := c.Ws.WriteMessage(1, []byte(fmt.Sprintf("%d", i))); err != nil {
                log.Println(err)
                return
            }
        }
        time.Sleep(1 * time.Second)
    }
    for _, c := range p.Clients {
        if err := c.Ws.WriteMessage(1, []byte("Start")); err != nil {
            log.Println(err)
            return
        }
        go p.readAndUpdateClientGameState(c, gameEnd)
    }
    go p.updateBall()
}

func(p *PrivateQueue) monitorEnd(gameEnd chan bool) {
    <- gameEnd
    log.Println("Game end ID=" + p.Id)
    for _, c := range p.Clients {
        c.Ws.WriteMessage(1, []byte("Game has ended"))
    }
}

var canvas = Canvas{width: 1400, height: 1000}
var ball = Ball{
    X: (canvas.width / 2) - 9,
    Y: (canvas.height / 2) - 9,
    DirectionX: IDLE,
    DirectionY: IDLE,
    Speed: 9,
}

/*
    if (this.ball.x - this.ball.width <= this.player.x && this.ball.x >= this.player.x - this.player.width) {
        if (this.ball.y <= this.player.y + this.player.height && this.ball.y + this.ball.height >= this.player.y) {
            this.ball.x = (this.player.x + this.ball.width);
            this.ball.moveX = DIRECTION.RIGHT;

            beep1.play();
        }
    }
*/

func (p *PrivateQueue) updateBall() {
    interval := time.Second / 60
    ball.Randomize(canvas.height)
    for {
        start := time.Now()
        switch ball.DirectionX {
        case LEFT:
            ball.X -= ball.Speed
        case RIGHT:
            ball.X += ball.Speed
        }
        switch ball.DirectionY {
        case UP:
            ball.Y -= int(float64(ball.Speed) / 1.5)
        case DOWN:
            ball.Y += int(float64(ball.Speed) / 1.5)
        }
        // handle ball out of bounds collisions
        if ball.X <= 0 { p._resetTurn() }
        if ball.X >= canvas.width - ball.Width() { p._resetTurn() }
        if ball.Y <= 0 { ball.DirectionY = DOWN }
        if ball.Y >= canvas.height - ball.Height() { ball.DirectionY = UP }
        time.Sleep(interval - time.Since(start))
    }
}

func (p *PrivateQueue) _resetTurn() {
    ball.Reset()
    ball.Randomize(canvas.height)

    time.Sleep(1 * time.Second)
}

func (p *PrivateQueue) readAndUpdateClientGameState(sender *Client, gameEnd chan bool) {
    for {
        var data GameDataReceive
        err := sender.Ws.ReadJSON(&data);
        if err != nil {
            handleError(err, gameEnd)
            return
        }

        payload := GameDataSend {
            PaddleX: data.PaddleX,
            BallY: ball.Y,
            BallX: ball.X,
        }

        for _, c := range p.Clients {
            if c == sender { continue }

            c.Ws.WriteJSON(payload)
        }
    }
}

func handleError(err error, gameEnd chan bool) {

    if err != nil {
    if ce, ok := err.(*websocket.CloseError);
        ok {
        switch ce.Code {
        case websocket.CloseNormalClosure:
        gameEnd <- true
        case websocket.CloseGoingAway:
            gameEnd <- true
        case websocket.CloseNoStatusReceived:
            gameEnd <- true
    }
        log.Println(err)
        return
    }
    }
}
