package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Craig-Turley/go-game-server/models"
	"github.com/gorilla/websocket"
)

var GlobalManager = models.NewPrivateGameManager()
var GlobalUpgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool { return true },
}

func homePageHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("HomePage"))
}

func createNewGame(w http.ResponseWriter, r *http.Request) {
    ws, err := GlobalUpgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        w.Write([]byte("Error connecting"))
        return
    }

    NewClient := models.Client {
        Ws:     ws,
        Status: models.INQUEUE,
    }
    id := GlobalManager.NewPrivateGame(&NewClient)

    rd := models.RoomData {
        ID: id,
        Error: "",
    }

    NewClient.Ws.WriteJSON(rd)

    log.Println("New game created with Id:", id);
}

func joinPrivateGame(w http.ResponseWriter, r *http.Request) {
    ws, err := GlobalUpgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        w.Write([]byte("Error Connecting"))
        return
    }

    id := strings.Split(r.URL.Path[1:], "/")[2]
    log.Println("Attempting to connect to room: ", id)

    GlobalManager.JoinGame(ws, id)
}

func setUpRoutes() {
    http.HandleFunc("/home", homePageHandler)
    http.HandleFunc("/ws/create_game", createNewGame)
    http.HandleFunc("/ws/join_game/", joinPrivateGame)
}

func main() {

    greeting := `
  ____        ____                      ____
 / ___| ___  / ___| __ _ _ __ ___   ___/ ___|  ___ _ ____   _____ _ __
| |  _ / _ \| |  _ / _' | '_ ' _ \ / _ \___ \ / _ \ '__\ \ / / _ \ '__|
| |_| | (_) | |_| | (_| | | | | | |  __/___) |  __/ |   \ V /  __/ |
 \____|\___/ \____|\__,_|_| |_| |_|\___|____/ \___|_|    \_/ \___|_|`

    setUpRoutes()

    fmt.Println(greeting)

    log.Fatal(http.ListenAndServe(":8080", nil))
}
