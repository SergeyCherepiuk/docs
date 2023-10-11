package http

import (
	"github.com/SergeyCherepiuk/docs/pkg/database/neo4j"
	"github.com/SergeyCherepiuk/docs/pkg/http/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Router struct{}

func (r Router) Build() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())

	userCreator := neo4j.NewUserCreator()
	userHandler := handlers.NewUserHandler(userCreator)

	v1 := e.Group("/api/v1")

	user := v1.Group("/user")
	user.POST("", userHandler.Create)

	return e
}
