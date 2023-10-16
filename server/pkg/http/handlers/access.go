package handlers

import (
	"context"
	"net/http"

	"github.com/SergeyCherepiuk/docs/domain"
	"github.com/SergeyCherepiuk/docs/pkg/http/handlers/internal"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type accessHandler struct {
	accessGranter domain.AccessGranter
	fileGetter    domain.FileGetter
}

func NewAccessHandler(
	accessGranter domain.AccessGranter,
	fileGetter domain.FileGetter,
) *accessHandler {
	return &accessHandler{
		accessGranter: accessGranter,
		fileGetter:    fileGetter,
	}
}

// TODO: Handle "updates", granting the different type of access to the same user
func (h accessHandler) Grant(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file id")
	}

	file, err := h.fileGetter.GetById(context.Background(), id.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, internal.ToSentence(err.Error()))
	}

	var access domain.Access
	if err := c.Bind(&access); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// TODO (?): Extract ownership checks into a middleware
	user, err := h.fileGetter.GetOwner(context.Background(), file)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
	}

	if user.Username != access.Granter {
		return echo.NewHTTPError(http.StatusInternalServerError, "User has no permissions to grant an access to the file")
	}

	// TODO: Validation

	if err := h.accessGranter.Grant(context.Background(), file, access); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
	}

	return c.NoContent(http.StatusCreated)
}
