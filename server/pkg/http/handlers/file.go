package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SergeyCherepiuk/docs/domain"
	"github.com/SergeyCherepiuk/docs/pkg/http/handlers/internal"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type fileHandler struct {
	fileCreator domain.FileCreator
	fileGetter  domain.FileGetter
	fileUpdater domain.FileUpdater
}

func NewFileHandler(
	fileCreator domain.FileCreator,
	fileGetter domain.FileGetter,
	fileUpdater domain.FileUpdater,
) *fileHandler {
	return &fileHandler{
		fileCreator: fileCreator,
		fileGetter:  fileGetter,
		fileUpdater: fileUpdater,
	}
}

func (h fileHandler) Create(c echo.Context) error {
	var file domain.File
	if err := c.Bind(&file); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	file.Id = uuid.NewString()

	fmt.Printf("%+v\n", file)

	// TODO: Validation

	if err := h.fileCreator.Create(context.Background(), file); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
	}

	return c.NoContent(http.StatusOK)
}

func (h fileHandler) Get(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file id")
	}

	file, err := h.fileGetter.GetById(context.Background(), id.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	return c.JSON(http.StatusOK, file)
}

type fileUpdates struct {
	NewName string `json:"newName"`
}

func (u fileUpdates) HasName() bool {
	return u.NewName != ""
}

func (h fileHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file id")
	}

	file, err := h.fileGetter.GetById(context.Background(), id.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	var updates fileUpdates
	if err := c.Bind(&updates); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if updates.HasName() {
		// TODO: Validation

		if err := h.fileUpdater.UpdateName(context.Background(), file, updates.NewName); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
		}
		return c.NoContent(http.StatusOK)
	}

	return c.NoContent(http.StatusBadRequest)
}
