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

	e.GET("/api/pointer", broadcast.Pointer)
	e.GET("/api/content", broadcast.Content)
	e.GET("/api/selection", broadcast.Selection)

	return e
}
