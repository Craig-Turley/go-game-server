package models

import (
	// "fmt"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type GameState string

const (
    INLOBBY  = "INLOBBY"
    PLAYING  = "PLAYING"
    FINISHED = "FINISHED"
)

type Queue interface {
    Add(cl *Client)
}

type PrivateQueue struct {
    Clients []*Client
    Id      string
    Host    *Client
    State   GameState
}

func (p *PrivateQueue) Add(cl *Client) {
    p.Clients = append(p.Clients, cl)
}

func (p *PrivateQueue) MonitorStart() {
    for {
        messageType, m, err := p.Host.Ws.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }

        log.Println(string(m))
        if string(m) == "GameStart" {
            for _, c := range p.Clients {
                if err := c.Ws.WriteMessage(messageType, []byte("GameStart")); err != nil {
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
}

func(p *PrivateQueue) monitorEnd(gameEnd chan bool) {
    <- gameEnd
    log.Println("Game end ID=" + p.Id)
    for _, c := range p.Clients {
        c.Ws.WriteMessage(1, []byte("Game has ended"))
    }
}

func (p *PrivateQueue) readAndUpdateClientGameState(sender *Client, gameEnd chan bool) {
    for {
        data, err := readSenderMessage(sender.Ws, p.State)
        if err != nil {
            handleError(err, gameEnd)
            return
        }

        for _, c := range p.Clients {
            if c == sender { continue }

            c.Ws.WriteJSON(data)
        }
    }
}

func readSenderMessage(ws *websocket.Conn, gameState GameState) (Payload, error) {
    switch gameState {
    case INLOBBY:
        var data GameStart
        err := ws.ReadJSON(&data)
        if err != nil {
            return nil, err
        }
        if data.Message == true {
            log.Println("Game Starting...")
        }
        return &data, nil
    case PLAYING:
        var data GameData
        err := ws.ReadJSON(&data)
        if err != nil {
            return nil, err
        }
        return &data, nil
    }
    var data GameData
    err := ws.ReadJSON(&data)
    if err != nil {
        return nil, err
    }
    return &data, nil
}

func handleError(err error, gameEnd chan bool) {

    if err != nil { if ce, ok := err.(*websocket.CloseError); ok { switch ce.Code { case websocket.CloseNormalClosure:
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
