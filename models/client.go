package models

import "github.com/gorilla/websocket"

type ClientStatus string

const (
    INQUEUE = "INQUEUE"
    //PLAYING = "PLAYING"
)

type Client struct {
    Ws     *websocket.Conn
    Status ClientStatus
}

func (c *Client) SetStatus(ns ClientStatus) {
    c.Status = ns
}
