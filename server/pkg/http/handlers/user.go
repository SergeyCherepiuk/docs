package handlers

import (
	"context"
	"net/http"

	"github.com/SergeyCherepiuk/docs/domain"
	"github.com/SergeyCherepiuk/docs/pkg/http/handlers/internal"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type userHandler struct {
	userCreator domain.UserCreator
	userGetter  domain.UserGetter
	userUpdater domain.UserUpdater
	userDeleter domain.UserDeleter
}

func NewUserHandler(
	userCreator domain.UserCreator,
	userGetter domain.UserGetter,
	userUpdater domain.UserUpdater,
	userDeleter domain.UserDeleter,
) *userHandler {
	return &userHandler{
		userCreator: userCreator,
		userGetter:  userGetter,
		userUpdater: userUpdater,
		userDeleter: userDeleter,
	}
}

func (h userHandler) Create(c echo.Context) error {
	var user domain.User
	if c.Bind(&user) != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// TODO: Validation

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash the password")
	}
	user.Password = string(hashedPassword)

	if err := h.userCreator.Create(context.Background(), user); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
	}

	return c.NoContent(http.StatusOK)
}

func (h userHandler) GetByUsername(c echo.Context) error {
	username := c.Param("username")

	user, err := h.userGetter.GetByUsername(context.Background(), username)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	return c.JSON(http.StatusOK, user)
}

type updates struct {
	NewUsername       string `json:"newUsername"`
	OldPassword       string `json:"oldPassword"`
	NewPassword       string `json:"newPassword"`
	NewPasswordRepeat string `json:"newPasswordRepeat"`
}

func (u updates) hasUsername() bool {
	return u.NewUsername != ""
}

func (u updates) hasPassword() bool {
	return u.OldPassword != "" && u.NewPassword != "" && u.NewPasswordRepeat != ""
}

func (h userHandler) Update(c echo.Context) error {
	username := c.Param("username")

	user, err := h.userGetter.GetByUsername(context.Background(), username)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User wasn't found")
	}

	var updates updates
	if c.Bind(&updates) != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if updates.hasUsername() {
		if err := h.userUpdater.UpdateUsername(context.Background(), user, updates.NewUsername); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
		}
		return c.NoContent(http.StatusOK)
	}

	if updates.hasPassword() {
		user, err := h.userGetter.GetByUsername(context.Background(), username)
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
		if err := h.userUpdater.UpdatePassword(context.Background(), user, string(hashedPassword)); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
		}
		return c.NoContent(http.StatusOK)
	}

	return c.NoContent(http.StatusBadRequest)
}

func (h userHandler) Delete(c echo.Context) error {
	username := c.Param("username")

	user, err := h.userGetter.GetByUsername(context.Background(), username)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "User wasn't found")
	}

	if err := h.userDeleter.Delete(context.Background(), user); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
	}

	return c.NoContent(http.StatusOK)
}
