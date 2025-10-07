package main

import (
	"net/http"

	"github.com/cldfn/wsbroadcast/app"
	"github.com/cldfn/wsbroadcast/server/routes"
	"go.uber.org/fx"
)

var GitCommit string = "dev"

func AsRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(app.RouteHandler)),
		fx.ResultTags(`group:"rhandlers"`),
	)
}

func main() {

	fx.New(
		fx.Supply(app.BuildInfo{
			GitCommit: GitCommit,
		}),
		fx.Provide(
			app.NewEnvProvider,
			app.NewHTTPServer,
			app.NewEnvConfig,
			app.NewBroadcaster,
			fx.Annotate(app.SetupRoutes, fx.ParamTags("", `group:"rhandlers"`)),
			AsRoute(routes.NewGlobalRoutes),
		),
		fx.Invoke(func(srv *http.Server) {

		}),
	).Run()

}
