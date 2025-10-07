package app

import (
	"github.com/cldfn/wsbroadcast/server"
	"github.com/google/uuid"
)

type Broadcaster struct {
	users         *server.LockedMap[uuid.UUID, server.WsClient]
	broadcastChan chan []byte
}

func (b *Broadcaster) PutUser(u *server.WsClient) {
	b.users.PutRef(u.Uid(), u)
}

// broadcast
func (b *Broadcaster) Broadcast(data []byte) {
	b.broadcastChan <- data
}

func NewBroadcaster() *Broadcaster {
	bcaster := &Broadcaster{}

	var usersInfo = server.NewLockedMap[uuid.UUID, server.WsClient]()

	broadcastChan := make(chan []byte, 10000)

	bcaster.users = usersInfo
	bcaster.broadcastChan = broadcastChan

	return bcaster
}
