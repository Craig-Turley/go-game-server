package models

type Payload interface {
    Payload()
}

type GameData struct {
    PaddleX int `json:"paddleX"`
}

func(g *GameData) Payload() {}

type GameStart struct {
    Message bool `json:"gameStart"`
}

func (g* GameStart) Payload() {}
