package main

import (
	"io"
	"log"
	"net/http"

	gin "github.com/cldfn/gina"
	"github.com/cldfn/utils"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func main() {

	var usersInfo = NewLockedMap[uuid.UUID, WsClient]()
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	broadcastChan := make(chan []byte, 10000)

	// go gobase.WithInterval(context.Background(), time.Millisecond*4000, func() {

	// 	clonedData := usersInfo.Copied()

	// 	tnow := time.Now()

	// 	for _, it := range clonedData {
	// 		it.conn.WriteControl(websocket.PingMessage, nil, tnow.Add(time.Second))
	// 	}

	// }).Run()

	utils.SafeGoroutine(func() {

		defer func() {
			log.Printf("broadcast routine finished")
		}()

		for dataToBroadcast := range broadcastChan {

			activeConns := usersInfo.Copied()

			for _, activeConn := range activeConns {
				if activeConn.LastError == nil {
					sendErr := activeConn.conn.WriteMessage(websocket.TextMessage, dataToBroadcast)
					if sendErr != nil {
						activeConn.LastError = sendErr
					}
				}
			}
		}
	})

	server := gin.New[WsServerContext]()
	server.SetContextDataIntializer(func(wsc *WsServerContext) {

	})

	server.GET("/ws", func(ctx *gin.Context[WsServerContext]) {

		uid := uuid.New()

		c, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)

		if err != nil {

			log.Printf(" -- ws upgrade err : %s", err.Error())

			ctx.JSON(500, gin.H{
				"msg": "unable to upgrade to ws conn",
				"err": err.Error(),
			})
			return
		}

		clientWrapp := WsClient{
			conn: c,
			uid:  uid,
		}

		usersInfo.Put(uid, clientWrapp)
	})

	server.POST("/broadcast", func(c *gin.Context[WsServerContext]) {

		data, readErr := io.ReadAll(c.Request.Body)

		if readErr != nil {
			broadcastChan <- data
		}

		c.JSON(200, gin.H{
			"msg": "sent",
		})
	})

	server.GET("/broadcast", func(c *gin.Context[WsServerContext]) {

		data := c.Query("msg")

		broadcastChan <- []byte(data)

		c.JSON(200, gin.H{
			"msg": "sent",
		})
	})

	server.GET("/send/:dest", func(c *gin.Context[WsServerContext]) {
		c.JSON(200, gin.H{
			"msg": "sent",
		})
	})

	nativeServer := &http.Server{
		Addr:    ":5599",
		Handler: server,
	}

	chainPath := "/etc/letsencrypt/live/ws.cldfn.com/fullchain.pem"
	keyPath := "/etc/letsencrypt/live/ws.cldfn.com/privkey.pem"

	if err := nativeServer.ListenAndServeTLS(chainPath, keyPath); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
