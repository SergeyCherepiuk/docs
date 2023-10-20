package handlers

import (
	"context"
	"net/http"

	"github.com/SergeyCherepiuk/docs/pkg/database/models"
	"github.com/SergeyCherepiuk/docs/pkg/database/neo4j"
	"github.com/SergeyCherepiuk/docs/pkg/http/handlers/internal"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AccessHandler struct{}

func (h AccessHandler) Grant(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file id")
	}

	file, err := neo4j.FileService.GetById(context.Background(), id.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, internal.ToSentence(err.Error()))
	}

	var accessBody models.Access
	if err := c.Bind(&accessBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	owner, err := neo4j.FileService.GetOwner(context.Background(), file)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
	}

	if owner.Username != accessBody.Granter {
		return echo.NewHTTPError(http.StatusInternalServerError, "User has no permissions to grant an access to the file")
	}

	// TODO: Validation

	receiver, err := neo4j.UserService.GetByUsername(context.Background(), accessBody.Receiver)
	if err != nil {
		return err
	}

	access, err := neo4j.AccessService.Get(context.Background(), file, receiver)
	if err != nil {
		if err := neo4j.AccessService.Grant(context.Background(), file, accessBody); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
		}
	} else {
		if err := neo4j.AccessService.UpdateLevel(context.Background(), file, access, accessBody.Level); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
		}
	}

	return c.NoContent(http.StatusCreated)
}

func (h AccessHandler) GetAccesses(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file id")
	}

	file, err := neo4j.FileService.GetById(context.Background(), id.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	accesses, err := neo4j.AccessService.GetAccesses(context.Background(), file)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
	}

	return c.JSON(http.StatusOK, accesses)
}

func (h AccessHandler) Revoke(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file id")
	}

	file, err := neo4j.FileService.GetById(context.Background(), id.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	username := c.Param("username")
	user, err := neo4j.UserService.GetByUsername(context.Background(), username)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	access, err := neo4j.AccessService.Get(context.Background(), file, user)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	if err := neo4j.AccessService.Revoke(context.Background(), file, access); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	return c.NoContent(http.StatusOK)
}
