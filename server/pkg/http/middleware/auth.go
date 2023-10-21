package middleware

import (
	"context"
	"net/http"

	"github.com/SergeyCherepiuk/docs/pkg/database/neo4j"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func NoSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session")
		if err != nil {
			return next(c)
		}

		id, err := uuid.Parse(cookie.Value)
		if err != nil {
			return next(c)
		}

		ctx := context.Background()
		sess := neo4j.NewSession(ctx)
		defer sess.Close(ctx)

		active, err := neo4j.SessionService.Check(ctx, sess, id)
		if !active || err != nil {
			return next(c)
		}

		// TODO: Set session owner's username to locals
		return echo.NewHTTPError(http.StatusBadRequest, "There is an active session")
	}
}
