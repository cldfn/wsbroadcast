package server

import (
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WsClient struct {
	conn *websocket.Conn
	uid  uuid.UUID

	LastPongReceived time.Time
	LastError        error
}

func (client *WsClient) Uid() uuid.UUID {
	return client.uid
}

func NewWsClient(conn *websocket.Conn, uid uuid.UUID) *WsClient {
	return &WsClient{
		conn:             conn,
		uid:              uid,
		LastPongReceived: time.Now(),
	}
}
