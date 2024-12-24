package main

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
