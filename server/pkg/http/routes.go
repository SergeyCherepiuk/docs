package http

import (
	"github.com/SergeyCherepiuk/docs/pkg/http/handlers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Router struct{}

func (r Router) Build() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())

	var (
		userHandler   = handlers.UserHandler{}
		fileHandler   = handlers.FileHandler{}
		accessHandler = handlers.AccessHandler{}
	)

	// TODO: Think about better API design
	v1 := e.Group("/api/v1")

	user := v1.Group("/user")
	user.POST("", userHandler.Create)
	user.GET("/:username", userHandler.GetByUsername)
	user.PUT("/:username", userHandler.Update)
	user.DELETE("/:username", userHandler.Delete)

	file := v1.Group("/file")
	file.POST("", fileHandler.Create)
	file.GET("/:id", fileHandler.Get)
	file.GET("/owner/:username", fileHandler.GetAll)
	file.PUT("/:id", fileHandler.Update)
	file.DELETE("/:id", fileHandler.Delete)
	file.DELETE("/owner/:username", fileHandler.DeleteAllForOwner)

	access := file.Group("/access")
	access.POST("/:id", accessHandler.Grant)
	access.GET("/:id", accessHandler.GetAccesses)
	access.DELETE("/:id/:username", accessHandler.Revoke)

	return e
}
