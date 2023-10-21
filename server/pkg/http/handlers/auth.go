package handlers

import (
	"context"
	"net/http"

	"github.com/SergeyCherepiuk/docs/pkg/database/models"
	"github.com/SergeyCherepiuk/docs/pkg/database/neo4j"
	"github.com/SergeyCherepiuk/docs/pkg/http/internal"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct{}

func (h AuthHandler) SignUp(c echo.Context) error {
	type RequestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var body RequestBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// TODO: Validation

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash the password")
	}

	ctx := context.Background()
	sess := neo4j.NewSession(ctx)
	defer sess.Close(ctx)

	tx, err := sess.BeginTransaction(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to start a database transaction")
	}

	user := models.User{
		Username: body.Username,
		Password: string(hashedPassword),
	}
	if err := neo4j.UserService.Create(ctx, tx, user); err != nil {
		tx.Rollback(ctx)
		return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
	}

	session := models.NewWeekSession(user.Username)
	if err := neo4j.SessionService.Create(ctx, tx, session); err != nil {
		tx.Rollback(ctx)
		return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
	}

	tx.Commit(ctx)
	c.SetCookie(&http.Cookie{
		Name:     "session",
		Value:    session.Id,
		Path:     "/",
		HttpOnly: true,
	})
	return c.JSON(http.StatusOK, session)
}

func (h AuthHandler) Login(c echo.Context) error {
	type RequestBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var body RequestBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	ctx := context.Background()
	sess := neo4j.NewSession(ctx)
	defer sess.Close(ctx)

	user, err := neo4j.UserService.GetByUsername(ctx, sess, body.Username)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Wrong password")
	}

	session := models.NewWeekSession(user.Username)
	if err := neo4j.SessionService.Create(ctx, sess, session); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
	}

	c.SetCookie(&http.Cookie{
		Name:     "session",
		Value:    session.Id,
		Path:     "/",
		HttpOnly: true,
	})
	return c.JSON(http.StatusOK, session)
}
