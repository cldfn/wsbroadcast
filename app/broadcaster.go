package app

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"time"

	"github.com/cldfn/wsbroadcast/server"
	"github.com/fatih/color"
	"github.com/google/uuid"
	"go.uber.org/fx"
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

	log.Printf("Broadcasting message of size %d bytes. addr is = %p", len(data), b.broadcastChan)

	b.broadcastChan <- data
}

type Worker struct {
	stop        chan struct{}
	queue       chan []byte
	broadcaster *Broadcaster
}

func (w *Worker) run() {

	errors := 0
	breakOnErrorsCount := 10

	for {

		if breakOnErrorsCount == errors {
			log.Printf("Too many errors occurred during broadcasting. Stopping...")
			break
		}

		func() {

			defer func() {

				rec := recover()

				if rec != nil {

					debug.PrintStack()

					log.Printf("[ERROR] while processing broadcast queue %v", rec)
					errors += 1
				}

			}()

			select {
			case <-w.stop:
				fmt.Println("Worker received stop signal.")
				errors = breakOnErrorsCount
				return
			case data := <-w.queue:

				copiedData := w.broadcaster.users.Copied()
				color.Green("copied users size : %d", len(copiedData))

				func() {
					for _, user := range copiedData {
						user.Write(data)
					}
				}()
			}
		}()

	}
}

func NewBroadcaster(lc fx.Lifecycle) *Broadcaster {
	bcaster := &Broadcaster{}

	var usersInfo = server.NewLockedMap[uuid.UUID, server.WsClient]()

	broadcastChan := make(chan []byte, 10000)

	bcaster.users = usersInfo
	bcaster.broadcastChan = broadcastChan

	w := &Worker{stop: make(chan struct{}), queue: bcaster.broadcastChan, broadcaster: bcaster}

	log.Printf(" -- bcaster chan addr : %p", bcaster.broadcastChan)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			log.Printf("Starting broadcaster")
			go w.run()

			return nil
		},
		OnStop: func(ctx context.Context) error {

			log.Printf("Stopping broadcaster")
			close(w.stop) // signal the goroutine to stop

			// Optionally wait for cleanup
			select {
			case <-time.After(5 * time.Second):
				fmt.Println("Worker stopped gracefully.")
			case <-ctx.Done():
				fmt.Println("Stop timed out.")
			}

			return nil
		},
	})

	return bcaster
}
