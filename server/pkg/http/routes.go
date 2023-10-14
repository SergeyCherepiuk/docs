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

	var (
		userCreator = neo4j.NewUserCreator()
		userGetter  = neo4j.NewUserGetter()
		userUpdater = neo4j.NewUserUpdater()
		userDeleter = neo4j.NewUserDeleter()
	)

	var (
		fileCreator = neo4j.NewFileCreator()
		fileGetter  = neo4j.NewFileGetter()
		fileUpdater = neo4j.NewFileUpdater()
	)

	var (
		userHandler = handlers.NewUserHandler(
			userCreator, userGetter, userUpdater, userDeleter,
		)
		fileHandler = handlers.NewFileHandler(
			fileCreator, fileGetter, fileUpdater,
		)
	)

	v1 := e.Group("/api/v1")

	user := v1.Group("/user")
	user.POST("", userHandler.Create)
	user.GET("/:username", userHandler.GetByUsername)
	user.PUT("/:username", userHandler.Update)
	user.DELETE("/:username", userHandler.Delete)

	file := v1.Group("/file")
	file.POST("", fileHandler.Create)
	file.GET("/:id", fileHandler.Get)
	file.PUT("/:id", fileHandler.Update)

	return e
}
