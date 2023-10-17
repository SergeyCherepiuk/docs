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
		fileDeleter = neo4j.NewFileDeleter()
	)

	var (
		accessGrater  = neo4j.NewAccessGranter()
		accessGetter  = neo4j.NewAccessGetter()
		accessUpdater = neo4j.NewAccessUpdater()
		accessRevoker = neo4j.NewAccessRevoker()
	)

	var (
		userHandler = handlers.NewUserHandler(
			userCreator, userGetter, userUpdater, userDeleter,
		)
		fileHandler = handlers.NewFileHandler(
			fileCreator, fileGetter, fileUpdater, fileDeleter, userGetter,
		)
		accessHandler = handlers.NewAccessHandler(
			accessGrater, accessGetter, accessUpdater, accessRevoker, fileGetter, userGetter,
		)
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
