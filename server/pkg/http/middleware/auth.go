package middleware

import (
	"context"
	"net/http"

	"github.com/SergeyCherepiuk/docs/pkg/database/neo4j"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	REQUIRE_SESSION    = 1
	REQUIRE_NO_SESSION = 2
)

func RequireSession() echo.MiddlewareFunc {
	return checkSession(REQUIRE_SESSION)
}

func RequireNoSession() echo.MiddlewareFunc {
	return checkSession(REQUIRE_NO_SESSION)
}

func checkSession(flag int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		unauthorized := func(c echo.Context) error {
			return c.NoContent(http.StatusUnauthorized)
		}

		var (
			onSessionPresent func(c echo.Context) error
			onSessionAbsent  func(c echo.Context) error
		)

		if flag == REQUIRE_SESSION {
			onSessionPresent = next
			onSessionAbsent = unauthorized
		} else if flag == REQUIRE_NO_SESSION {
			onSessionPresent = unauthorized
			onSessionAbsent = next
		}

		return func(c echo.Context) error {
			cookie, err := c.Cookie("session")
			if err != nil {
				return onSessionAbsent(c)
			}

			id, err := uuid.Parse(cookie.Value)
			if err != nil {
				return onSessionAbsent(c)
			}

			ctx := context.Background()
			sess := neo4j.NewSession(ctx)
			defer sess.Close(ctx)

			user, err := neo4j.SessionService.Check(ctx, sess, id)
			if err != nil {
				return onSessionAbsent(c)
			}

			c.Set("user", user)
			return onSessionPresent(c)
		}
	}
}
