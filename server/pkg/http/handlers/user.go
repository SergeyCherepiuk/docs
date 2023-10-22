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

type UserHandler struct{}

func (h UserHandler) GetByUsername(c echo.Context) error {
	username := c.Param("username")

	ctx := context.Background()
	sess := neo4j.NewSession(ctx)
	defer sess.Close(ctx)

	user, err := neo4j.UserService.GetByUsername(ctx, sess, username)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	return c.JSON(http.StatusOK, user)
}

type userUpdates struct {
	NewUsername       string `json:"newUsername"`
	OldPassword       string `json:"oldPassword"`
	NewPassword       string `json:"newPassword"`
	NewPasswordRepeat string `json:"newPasswordRepeat"`
}

func (u userUpdates) hasUsername() bool {
	return u.NewUsername != ""
}

func (u userUpdates) hasPassword() bool {
	return u.OldPassword != "" && u.NewPassword != "" && u.NewPasswordRepeat != ""
}

func (h UserHandler) Update(c echo.Context) error {
	user, ok := c.Get("user").(models.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "User wasn't found")
	}

	ctx := context.Background()
	sess := neo4j.NewSession(ctx)
	defer sess.Close(ctx)

	var updates userUpdates
	if c.Bind(&updates) != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if updates.hasUsername() {
		// TODO: Validation

		if err := neo4j.UserService.UpdateUsername(ctx, sess, user, updates.NewUsername); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
		}
		return c.NoContent(http.StatusOK)
	}

	if updates.hasPassword() {
		// TODO: Validation

		user, err := neo4j.UserService.GetByUsername(ctx, sess, user.Username)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(updates.OldPassword)); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Wrong password")
		}

		if updates.NewPassword != updates.NewPasswordRepeat {
			return echo.NewHTTPError(http.StatusBadRequest, "New passwords aren't the same")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updates.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash the password")
		}
		if err := neo4j.UserService.UpdatePassword(ctx, sess, user, string(hashedPassword)); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
		}
		return c.NoContent(http.StatusOK)
	}

	return c.NoContent(http.StatusBadRequest)
}

func (h UserHandler) Delete(c echo.Context) error {
	user, ok := c.Get("user").(models.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "User wasn't found")
	}

	ctx := context.Background()
	sess := neo4j.NewSession(ctx)
	defer sess.Close(ctx)

	if err := neo4j.UserService.Delete(ctx, sess, user); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
	}

	return c.NoContent(http.StatusOK)
}
