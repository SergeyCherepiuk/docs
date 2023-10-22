package http

import (
	"github.com/SergeyCherepiuk/docs/pkg/http/handlers"
	"github.com/SergeyCherepiuk/docs/pkg/http/middleware"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

type Router struct{}

func (r Router) Build() *echo.Echo {
	e := echo.New()
	e.Use(echomiddleware.Logger())

	var (
		authHandler   = handlers.AuthHandler{}
		userHandler   = handlers.UserHandler{}
		fileHandler   = handlers.FileHandler{}
		accessHandler = handlers.AccessHandler{}
	)

	v1 := e.Group("/api/v1")

	auth := v1.Group("/auth")
	auth.Use(middleware.RequireNoSession())
	auth.POST("/signup", authHandler.SignUp)
	auth.POST("/login", authHandler.Login)

	v1.Use(middleware.RequireSession())

	v1.POST("/auth/logout", authHandler.LogOut)

	user := v1.Group("/user")
	user.GET("/:username", userHandler.GetByUsername)
	user.PUT("", userHandler.Update)
	user.DELETE("", userHandler.Delete)

	file := v1.Group("/files")
	file.POST("", fileHandler.Create)
	file.GET("/:id", fileHandler.Get, middleware.RequireAtLeastRAccess)
	file.GET("", fileHandler.GetAll)
	file.PUT("/:id", fileHandler.Update, middleware.RequireAtLeastRWAccess)
	file.DELETE("/:id", fileHandler.Delete, middleware.RequireOwnerAccess)

	access := file.Group("/access")
	access.POST("/:id", accessHandler.Grant, middleware.RequireOwnerAccess)
	access.GET("/:id", accessHandler.GetAccesses, middleware.RequireAtLeastRAccess)
	access.DELETE("/:id/:username", accessHandler.Revoke, middleware.RequireOwnerAccess)

	return e
}
