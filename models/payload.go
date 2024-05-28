package models

import (
	"math/rand"
)

type Payload interface {
    Payload()
}

type GameDataReceive struct {
    PaddleY float32 `json:"paddleY"`
    PaddleX float32 `json:"PaddleX"`
}

func(g *GameDataReceive) Payload() {}

type GameDataSend struct {
    PaddleY        float32 `json:"paddleY"`
    BallX          float32 `json:"ballX"`
    BallY          float32 `json:"ballY"`
    PlayerOneScore int     `json:"playerOneScore"`
    PlayerTwoScore int     `json:"playerTwoScore"`
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
    X          float32   `json:"ballX"`
    Y          float32   `json:"ballY"`
    DirectionX DIRECTION `json:"directionX"`
    DirectionY DIRECTION `json:"directionY"`
    Speed      float32   `json:"speed"`
}

func (b *Ball) Payload()    {}
func (b *Ball) Height() float32 { return 18 }
func (b *Ball) Width()  float32 { return 18 }
func (b *Ball) Reset() {
    b.X = (canvas.width / 2) - 9
    b.Y = (canvas.height / 2) - 9
    b.DirectionX = IDLE
    b.DirectionY = IDLE
    b.Speed = 9
}
func (b *Ball) Randomize(height float32) {
    random := []DIRECTION{UP, DOWN, LEFT, RIGHT}
    // b.moveX = add turn based directions
    b.DirectionX = random[rand.Intn(4-2) + 2]
    b.DirectionY = random[rand.Intn(2)]
    b.Y = rand.Float32() * (height - 200) + 200
}

type Score struct {
    PlayerOne   int  `json:"playerOne"`
    PlayerTwo   int  `json:"playerTwo"`
}

func (s *Score) Payload() {}
