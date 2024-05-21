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
    NewQueue := PrivateQueue {
        Clients: []*Client{host},
        Id:      id,
        Host:    host,
        State:   INLOBBY,
    }

    pgm.Games[id] = &NewQueue
    go NewQueue.MonitorStart();
    log.Println("Waiting to start game with id: " + id)

    return id

}

func (pgm *PrivateGameManager) findGame(id string) (*PrivateQueue, error) {
    room, ok := pgm.Games[id]
    if !ok {
        log.Println("Error joining game with ID= " + id)
        return nil, errors.New("Error finding game with ID=" + id)
    }
    return room, nil
}

func (pgm *PrivateGameManager) JoinGame(ws *websocket.Conn, id string) {
    room, err := pgm.findGame(id)
    if err != nil {
        log.Println(err)
        ws.WriteMessage(1, []byte(err.Error()))
        return
    }

    cl := Client {
        Ws:     ws,
        Status: INQUEUE,
    }

    room.Add(&cl)
    log.Println("Successfully joining game with ID= " + id)
    cl.Ws.WriteMessage(1, []byte("Successfully joining game with ID= " + id))
}

func randomId() string {
    return fmt.Sprintf("%d", 1000 + rand.Intn(9999-1000))
}
