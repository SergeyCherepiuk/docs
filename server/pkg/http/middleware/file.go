package middleware

import (
	"context"
	"net/http"

	"github.com/SergeyCherepiuk/docs/pkg/database/models"
	"github.com/SergeyCherepiuk/docs/pkg/database/neo4j"
	"github.com/SergeyCherepiuk/docs/pkg/http/internal"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const ownerAccess = "O"

func RequireAtLeastRAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		level, err := getAccessLevel(c)
		if err != nil {
			return err
		}

		if !isAtLeastRAccess(level) {
			return echo.NewHTTPError(http.StatusUnauthorized, "At least 'read' access required")
		}

		return next(c)
	}
}

func RequireAtLeastRWAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		level, err := getAccessLevel(c)
		if err != nil {
			return err
		}

		if !isAtLeastRWAccess(level) {
			return echo.NewHTTPError(http.StatusUnauthorized, "At least 'read&write' access required")
		}

		return next(c)
	}
}

func RequireOwnerAccess(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		level, err := getAccessLevel(c)
		if err != nil {
			return err
		}

		if !isOwnerAccess(level) {
			return echo.NewHTTPError(http.StatusUnauthorized, "Owner access required")
		}

		return next(c)
	}
}

func isAtLeastRAccess(level string) bool {
	return level == models.RAcess || isAtLeastRWAccess(level)
}

func isAtLeastRWAccess(level string) bool {
	return level == models.RWAccess || isOwnerAccess(level)
}

func isOwnerAccess(level string) bool {
	return level == ownerAccess
}

func getAccessLevel(c echo.Context) (string, error) {
	user, ok := c.Get("user").(models.User)
	if !ok {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "User wasn't found")
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return "", echo.NewHTTPError(http.StatusBadRequest, "Invalid file id")
	}

	ctx := context.Background()
	sess := neo4j.NewSession(ctx)
	defer sess.Close(ctx)

	file, err := neo4j.FileService.GetById(ctx, sess, id)
	if err != nil {
		return "", echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	owner, err := neo4j.FileService.GetOwner(ctx, sess, file)
	if err == nil && owner.Username == user.Username {
		return ownerAccess, nil
	}

	access, err := neo4j.AccessService.Get(ctx, sess, file, user)
	if err != nil {
		return "", echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	return access.Level, nil
}
