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
	accessGetter  domain.AccessGetter
	accessUpdater domain.AccessUpdater
	fileGetter    domain.FileGetter
	userGetter domain.UserGetter
}

func NewAccessHandler(
	accessGranter domain.AccessGranter,
	accessGetter domain.AccessGetter,
	accessUpdater domain.AccessUpdater,
	fileGetter domain.FileGetter,
	userGetter domain.UserGetter,
) *accessHandler {
	return &accessHandler{
		accessGranter: accessGranter,
		accessGetter:  accessGetter,
		accessUpdater: accessUpdater,
		fileGetter:    fileGetter,
		userGetter: userGetter,
	}
}

func (h accessHandler) Grant(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file id")
	}

	file, err := h.fileGetter.GetById(context.Background(), id.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, internal.ToSentence(err.Error()))
	}

	var accessBody domain.Access
	if err := c.Bind(&accessBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	owner, err := h.fileGetter.GetOwner(context.Background(), file)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
	}

	if owner.Username != accessBody.Granter {
		return echo.NewHTTPError(http.StatusInternalServerError, "User has no permissions to grant an access to the file")
	}

	// TODO: Validation

	receiver, err := h.userGetter.GetByUsername(context.Background(), accessBody.Receiver)
	if err != nil {
		return err
	}

	access, err := h.accessGetter.Get(context.Background(), file, receiver)
	if err != nil {	
		if err := h.accessGranter.Grant(context.Background(), file, accessBody); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
		}
	} else {
		if err := h.accessUpdater.UpdateLevel(context.Background(), file, access, accessBody.Level); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
		}
	}

	return c.NoContent(http.StatusCreated)
}

func (h accessHandler) GetAccesses(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file id")
	}

	file, err := h.fileGetter.GetById(context.Background(), id.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	accesses, err := h.accessGetter.GetAccesses(context.Background(), file)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
	}

	return c.JSON(http.StatusOK, accesses)
}
