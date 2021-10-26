package main

import (
	"github.com/rasimogluali/SLAprover/internal/config"
	"github.com/rasimogluali/SLAprover/internal/server"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			config.NewConfig,
			server.NewServer,
			server.NewSLAproverHandler,
		),
		fx.Invoke(
			server.RegisterRoutes,
		),
	)
	app.Run()

}
