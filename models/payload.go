package models

import (
	"math/rand"
)

type Payload interface {
    Payload()
}

type GameDataReceive struct {
    PaddleX    float32 `json:"paddleY"`
}

func(g *GameDataReceive) Payload() {}

type GameDataSend struct {
    PaddleX    float32 `json:"paddleY"`
    BallX      int     `json:"ballX"`
    BallY      int     `json:"ballY"`
}
func(g *GameDataSend) Payload() {}

type GameStart struct {
    Message   bool `json:"gameStart"`
    Player    int  `json:"player"`
}

func (g* GameStart) Payload() {}

type RoomData struct {
    ID    string `json:"roomId"`
    Error string `json:"error"`
}

func (r *RoomData) Payload() {}

type DIRECTION int

const (
    IDLE = iota
    UP
    DOWN
    LEFT
    RIGHT
)

type Ball struct {
    X          int       `json:"ballX"`
    Y          int       `json:"ballY"`
    DirectionX DIRECTION `json:"directionX"`
    DirectionY DIRECTION `json:"directionY"`
    Speed      int       `json:"speed"`
}

func (b *Ball) Payload()    {}
func (b *Ball) Height() int { return 18 }
func (b *Ball) Width()  int { return 18 }
func (b *Ball) Reset() {
    b.X = (canvas.width / 2) - 9
    b.Y = (canvas.height / 2) - 9
    b.DirectionX = IDLE
    b.DirectionY = IDLE
    b.Speed = 9
}
func (b *Ball) Randomize(height int) {
    random := []DIRECTION{UP, DOWN, LEFT, RIGHT}
    // b.moveX = add turn based directions
    b.DirectionX = random[rand.Intn(4-2) + 2]
    b.DirectionY = random[rand.Intn(2)]
    ball.Y = rand.Intn(height - 200) + 200
}
