package app

import (
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(bi BuildInfo, rhandlers []RouteHandler, conf *EnvConfig) http.Handler {

	ginRouter := gin.New()

	ginRouter.Use(gin.Logger())
	ginRouter.Use(gin.CustomRecovery(func(ctx *gin.Context, err any) {

		if err != nil {

			// if appContext.PanicStack {

			stack := debug.Stack()

			ctx.JSON(500, gin.H{
				"panic": err,
				"stack": stack,
			})
			// } else {
			// 	ctx.JSON(500, gin.H{
			// 		"msg": "internal error",
			// 	})
			// }
		}
	}))

	// cors
	ginRouter.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	ginRouter.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"version": bi.GitCommit,
		})
	})

	for _, it := range rhandlers {
		it.Handler(ginRouter.Group(it.Path()))
	}

	return ginRouter
}
