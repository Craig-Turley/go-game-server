package models

import "github.com/gorilla/websocket"

type ClientStatus string

const (
    INQUEUE = "INQUEUE"
    //PLAYING = "PLAYING"
)

type Client struct {
    Ws             *websocket.Conn
    Status         ClientStatus
    Paddle         PongPaddle
    PlayerPosition int
}

func (c *Client) SetStatus(ns ClientStatus) {
    c.Status = ns
}

func (c *Client) SetPosition(position int) {
    c.PlayerPosition = position
}

type PongPaddle struct {
    X float32
    Y float32
}

func (p *PongPaddle) Set(x float32, y float32) {
    p.X = x
    p.Y = y
}

func (p *PongPaddle) Height() float32 { return 70 }
func (p *PongPaddle) Width()  float32 { return 18 }
