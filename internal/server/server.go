package server

import (
	"context"
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rasimogluali/SLAprover/internal/config"
	"go.uber.org/fx"
)

func NewServer(lc fx.Lifecycle, conf *config.Config) *echo.Echo {
	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "[${time}] ${method} ${host}${uri} status=${status}\n",
	}))

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Print(fmt.Sprintf("Listening at port %s...", conf.App.Port))
			go e.Start(fmt.Sprintf(":%s", conf.App.Port))
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Print("Server is shutting down")
			return e.Shutdown(ctx)
		},
	})

	return e
}
