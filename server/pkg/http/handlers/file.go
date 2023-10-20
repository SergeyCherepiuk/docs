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

type FileHandler struct{}

func (h FileHandler) Create(c echo.Context) error {
	// TODO: Replace hardcoded username with getting it from the sessionId
	username := "johndoe"

	user, err := neo4j.UserService.GetByUsername(context.Background(), username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, internal.ToSentence(err.Error()))
	}

	var file models.File
	if err := c.Bind(&file); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	file.Id = uuid.NewString()

	// TODO: Validation

	if err := neo4j.FileService.Create(context.Background(), file, user); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
	}

	return c.NoContent(http.StatusCreated)
}

func (h FileHandler) Get(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file id")
	}

	file, err := neo4j.FileService.GetById(context.Background(), id.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	user, err := neo4j.FileService.GetOwner(context.Background(), file)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	response := struct {
		File  models.File `json:"file"`
		Owner models.User `json:"owner"`
	}{file, user}
	return c.JSON(http.StatusOK, response)
}

func (h FileHandler) GetAll(c echo.Context) error {
	username := c.Param("username")

	user, err := neo4j.UserService.GetByUsername(context.Background(), username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, internal.ToSentence(err.Error()))
	}

	files, err := neo4j.FileService.GetAllForOwner(context.Background(), user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, internal.ToSentence(err.Error()))
	}

	response := struct {
		Owner models.User   `json:"owner"`
		Files []models.File `json:"files"`
	}{user, files}
	return c.JSON(http.StatusOK, response)
}

type fileUpdates struct {
	NewName string `json:"newName"`
}

func (u fileUpdates) HasName() bool {
	return u.NewName != ""
}

func (h FileHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file id")
	}

	file, err := neo4j.FileService.GetById(context.Background(), id.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	var updates fileUpdates
	if err := c.Bind(&updates); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if updates.HasName() {
		// TODO: Validation

		if err := neo4j.FileService.UpdateName(context.Background(), file, updates.NewName); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
		}
		return c.NoContent(http.StatusOK)
	}

	return c.NoContent(http.StatusBadRequest)
}

func (h FileHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file's id")
	}

	file, err := neo4j.FileService.GetById(context.Background(), id.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	if err := neo4j.FileService.Delete(context.Background(), file); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
	}

	return c.NoContent(http.StatusOK)
}

func (h FileHandler) DeleteAllForOwner(c echo.Context) error {
	username := c.Param("username")

	user, err := neo4j.UserService.GetByUsername(context.Background(), username)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	if err := neo4j.FileService.DeleteAllForOwner(context.Background(), user); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
	}

	return c.NoContent(http.StatusOK)
}
