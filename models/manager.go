package models

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"

	"github.com/gorilla/websocket"
)

type PrivateGameManager struct {
    Games map[string]*PrivateQueue
    mu    sync.Mutex
}

func NewPrivateGameManager() *PrivateGameManager {
    return &PrivateGameManager{
        Games: make(map[string]*PrivateQueue),
    }
}

func (pgm *PrivateGameManager) NewPrivateGame(host *Client) string {
    pgm.mu.Lock()
    defer pgm.mu.Unlock()

    id := randomId()
    ball := Ball{
        X: (canvas.width / 2) - 9,
        Y: (canvas.height / 2) - 9,
        DirectionX: IDLE,
        DirectionY: IDLE,
        Speed: 9,
    }
    NewQueue := PrivateQueue {
        Manager: pgm,
        Clients: []*Client{host},
        Id:      id,
        Host:    host,
        State:   INLOBBY,
        ball:    ball,
    }

    pgm.Games[id] = &NewQueue
    go NewQueue.MonitorStart();
    log.Println("Waiting to start game with id: " + id)

    return id

}

func (pgm *PrivateGameManager) findGame(id string) (*PrivateQueue, error) {
    room, ok := pgm.Games[id]
    if !ok {
        return nil, errors.New("Error finding game with ID: " + id)
    }
    return room, nil
}

func (pgm *PrivateGameManager) JoinGame(ws *websocket.Conn, id string) {
    room, err := pgm.findGame(id)
    if err != nil {
        log.Println(err)
        rd := RoomData {
            ID: "",
            Error: err.Error(),
        }
        ws.WriteJSON(rd)
        return
    }

    cl := Client {
        Ws:     ws,
        Status: INQUEUE,
    }

    rd := RoomData {
        ID: id,
        Error: "",
    }

    room.Add(&cl)
    log.Println("Successfully joining game with ID: " + id)
    cl.Ws.WriteJSON(rd)
}

func (pgm *PrivateGameManager) CleanUp(id string) {
    pgm.mu.Lock()
    delete(pgm.Games, id)
    pgm.mu.Unlock()
}

func randomId() string {
    return fmt.Sprintf("%d", 1000 + rand.Intn(9999-1000))
}
