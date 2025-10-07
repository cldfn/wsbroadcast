package routes

import (
	"io"
	"log"
	"net/http"

	"github.com/cldfn/wsbroadcast/app"
	"github.com/cldfn/wsbroadcast/server"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type GlobalRoutes struct {
	broadcaster *app.Broadcaster
}

var (
	_ app.RouteHandler = &GlobalRoutes{}
)

// Handler implements app.RouteHandler.
func (g *GlobalRoutes) Handler(router gin.IRouter) {

	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	router.GET("/connect", func(ctx *gin.Context) {

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

		clientWrapp := server.NewWsClient(
			c,
			uid,
		)

		g.broadcaster.PutUser(clientWrapp)
	})

	router.POST("/broadcast", func(c *gin.Context) {

		data, readErr := io.ReadAll(c.Request.Body)

		if readErr != nil {
			g.broadcaster.Broadcast(data)
		}

		c.JSON(200, gin.H{
			"msg": "sent",
		})
	})

	router.GET("/broadcast", func(c *gin.Context) {

		data := c.Query("msg")

		g.broadcaster.Broadcast([]byte(data))

		c.JSON(200, gin.H{
			"msg": "sent",
		})
	})

	router.GET("/send/:dest", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "sent",
		})
	})

}

// Path implements app.RouteHandler.
func (g *GlobalRoutes) Path() string {
	return "/"
}

func NewGlobalRoutes() *GlobalRoutes {
	r := new(GlobalRoutes)

	return r
}
