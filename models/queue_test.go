package models

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

var GlobalManagerQueueTest = NewPrivateGameManager()

func TestMonitorStartAndStartGame(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(monitorStartHandler))
    defer ts.Close()

    u := "ws" + strings.TrimPrefix(ts.URL, "http")
    ws, _ ,err := websocket.DefaultDialer.Dial(u, nil)
    if err != nil {
        t.Fatalf("%v", err)
    }
    defer ws.Close()

    var roomDataResponse RoomData
    if err := ws.ReadJSON(&roomDataResponse); err != nil {
        t.Fatalf("Error Reading RoomData JSON from server")
    }

    defer GlobalManager.CleanUp(roomDataResponse.ID)

    if err := ws.WriteMessage(1, []byte("GameStart")); err != nil {
        t.Fatalf("Error writing GameStart to server")
    }

    var gameStartResponse GameStart
    if err := ws.ReadJSON(&gameStartResponse); err != nil {
        t.Fatalf("Error reading Start Reponse from server")
    }

    if !gameStartResponse.Message {
        t.Fatalf("Game start message not recieved from server")
    }

    for i := 3; i > 0; i-- {
        _, p, err := ws.ReadMessage()
        if err != nil {
            t.Fatalf("Error reading countdown from server")
        }

        if string(p) != fmt.Sprintf("%d", i) {
            t.Fatalf("Countdown not matching expected number. Wanted %s got %s", fmt.Sprintf("%d", i), string(p))
        }
    }

    _, p, err := ws.ReadMessage()
    if err != nil {
        t.Fatalf("Error reading GameStart message from server")
    }

    if string(p) != "Start" {
        t.Fatalf("Unexpected response from the server. Wanted %s got %s", "Start", string(p))
    }

}

func monitorStartHandler(w http.ResponseWriter, r *http.Request) {
    upgrader := websocket.Upgrader{}

    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Fatalf(err.Error())
    }

    id := GlobalManagerQueueTest.NewPrivateGame(&Client{
        Ws: ws,
        Status: INQUEUE,
    })

    ws.WriteJSON(&RoomData{
        ID: id,
        Error: "",
    })

}
/*
func TestReadAndUpdateClientGameState(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(monitorStartHandler))
    defer ts.Close()

    u := "ws" + strings.TrimPrefix(ts.URL, "http")
    host, _ ,err := websocket.DefaultDialer.Dial(u, nil)
    if err != nil {
        t.Fatalf("%v", err)
    }
    defer host.Close()

    var roomDataResponse RoomData
    defer GlobalManager.CleanUp(roomDataResponse.ID)

    u = "ws" + strings.TrimPrefix(ts.URL, "http") + fmt.Sprintf("/join_game/%s", roomDataResponse.ID)
    ws, _, err := websocket.DefaultDialer.Dial(u, nil)

    host.WriteMessage(1, []byte("GameStart"))

    var gameStartResponse GameStart
    host.ReadJSON(&gameStartResponse)

    for i := 4; i > 0; i-- {
        host.ReadMessage()
        ws.ReadMessage()
    }

    for i := 0; i < 3; i++ {
        ws.WriteMessage(1, []byte(fmt.Sprintf("%d", i)))
        _, p, err := host.ReadMessage()
        if err != nil {
            t.Fatalf(err.Error())
        }
        if string(p) != fmt.Sprintf("%d", i) {
            t.Fatalf("Unexpected response received from server")
        }
    }


}
*/


