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
}

func NewFileHandler(
	fileCreator domain.FileCreator,
	fileGetter domain.FileGetter,
) *fileHandler {
	return &fileHandler{
		fileCreator: fileCreator,
		fileGetter:  fileGetter,
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
