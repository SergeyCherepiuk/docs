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
}

func NewUserHandler(userCreator domain.UserCreator) *userHandler {
	return &userHandler{userCreator: userCreator}
}

func (h userHandler) Create(c echo.Context) error {
	var user domain.User
	if c.Bind(&user) != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

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
