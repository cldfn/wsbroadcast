package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type RouteHandler interface {
	Handler(g gin.IRouter)
	Path() string
}

func NewHTTPServer(lc fx.Lifecycle, routes http.Handler, envConfig *EnvConfig) *http.Server {

	portValue := envConfig.ApiPort

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", portValue),
		Handler: routes,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			log.Println("Starting HTTP server at", srv.Addr)
			go srv.Serve(ln)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return srv
}
