package http

import (
	"github.com/SergeyCherepiuk/docs/pkg/http/broadcast"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Router struct{}

func (r Router) Build() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())

	e.GET("/api/connect", broadcast.Connect)

	return e
}
