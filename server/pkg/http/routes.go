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
	userGetter := neo4j.NewUserGetter()
	userChecker := neo4j.NewUserChecker()
	userUpdater := neo4j.NewUserUpdater()
	userDeleter := neo4j.NewUserDeleter()

	userHandler := handlers.NewUserHandler(
		userCreator, userGetter, userChecker, userUpdater, userDeleter,
	)

	v1 := e.Group("/api/v1")

	user := v1.Group("/user")
	user.GET("/:username", userHandler.GetByUsername)
	user.POST("", userHandler.Create)
	user.PUT("/:username", userHandler.Update)
	user.DELETE("/:username", userHandler.Delete)

	return e
}
