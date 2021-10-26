package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func Register(h *SLAproverHandler, group *echo.Group) {
	group.POST("/check-registry", h.checkRegistry)
}

func RegisterRoutes(engine *echo.Echo, h *SLAproverHandler) {
	engine.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})
	Register(h, engine.Group(""))
}
