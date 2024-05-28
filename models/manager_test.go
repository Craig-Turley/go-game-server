package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

var GlobalManager = NewPrivateGameManager()

func TestNewPrivateGame(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(createNewGameHandler))
    defer ts.Close()

    u := "ws" + strings.TrimPrefix(ts.URL, "http")
    ws, _ ,err := websocket.DefaultDialer.Dial(u, nil)
    if err != nil {
        t.Fatalf("%v", err)
    }
    defer ws.Close()

    _, p, err := ws.ReadMessage()
    if err != nil {
        t.Fatalf("%v", err)
    }

    var response RoomData
    json.Unmarshal(p, &response)

    if response.ID == "" {
        t.Fatalf("Game Id not returned")
    }

    _, err = GlobalManager.findGame(response.ID)
    if err != nil {
        t.Fatalf(err.Error())
    }

    GlobalManager.CleanUp(response.ID)
    GlobalManager = NewPrivateGameManager()
}

func TestJoinGame(t *testing.T) {
    ts := httptest.NewServer(http.HandlerFunc(createNewGameHandler))

    u := "ws" + strings.TrimPrefix(ts.URL, "http")
    fmt.Println(u)
    host, _ , _ := websocket.DefaultDialer.Dial(u, nil)
    defer host.Close()

    _, p, _ := host.ReadMessage()

    var response RoomData
    json.Unmarshal(p, &response)
    ts.Close()

    ts = httptest.NewServer(http.HandlerFunc(joinNewGameHandler))
    u = "ws" + strings.TrimPrefix(ts.URL, "http") + fmt.Sprintf("/join_game/%s", response.ID)
    ws, _, err := websocket.DefaultDialer.Dial(u, nil)
    if err != nil {
        t.Fatalf("%v", err)
    }

    var joinResponse RoomData
    ws.ReadJSON(&joinResponse)
    if joinResponse.Error != "" {
        t.Fatalf("Error when joining room: %s", joinResponse.Error)
    }

    if joinResponse.ID == "" || joinResponse.ID != response.ID {
        t.Fatalf("RoomId incorrect")
    }

}

func upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
    upgrader := websocket.Upgrader{}
    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return nil, err
    }

    return ws, nil

}

func createNewGameHandler(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrade(w, r)
    if err != nil {
        return
    }

    NewClient := Client {
        Ws:     ws,
        Status: INQUEUE,
    }
    id := GlobalManager.NewPrivateGame(&NewClient)

    rd := RoomData {
        ID: id,
        Error: "",
    }

    ws.WriteJSON(rd)
}

func joinNewGameHandler(w http.ResponseWriter, r *http.Request) {
    ws, err := upgrade(w, r)
    if err != nil {
        return
    }

    id := strings.Split(r.URL.Path[1:], "/")[1]
    GlobalManager.JoinGame(ws, id)
}
