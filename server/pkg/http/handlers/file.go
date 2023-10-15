package handlers

import (
	"context"
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
	fileDeleter domain.FileDeleter
	userGetter  domain.UserGetter
}

func NewFileHandler(
	fileCreator domain.FileCreator,
	fileGetter domain.FileGetter,
	fileUpdater domain.FileUpdater,
	fileDeleter domain.FileDeleter,
	userGetter domain.UserGetter,
) *fileHandler {
	return &fileHandler{
		fileCreator: fileCreator,
		fileGetter:  fileGetter,
		fileUpdater: fileUpdater,
		fileDeleter: fileDeleter,
		userGetter:  userGetter,
	}
}

func (h fileHandler) Create(c echo.Context) error {
	// TODO: Replace hardcoded username with getting it from the sessionId
	username := "johndoe"

	user, err := h.userGetter.GetByUsername(context.Background(), username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, internal.ToSentence(err.Error()))
	}

	var file domain.File
	if err := c.Bind(&file); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}
	file.Id = uuid.NewString()

	// TODO: Validation

	if err := h.fileCreator.Create(context.Background(), file, user); err != nil {
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

	user, err := h.fileGetter.GetOwner(context.Background(), file)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	response := struct {
		File  domain.File `json:"file"`
		Owner domain.User `json:"owner"`
	}{file, user}
	return c.JSON(http.StatusOK, response)
}

func (h fileHandler) GetAll(c echo.Context) error {
	username := c.Param("username")

	user, err := h.userGetter.GetByUsername(context.Background(), username)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, internal.ToSentence(err.Error()))
	}

	files, err := h.fileGetter.GetAllForOwner(context.Background(), user)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, internal.ToSentence(err.Error()))
	}

	response := struct {
		Owner domain.User   `json:"owner"`
		Files []domain.File `json:"files"`
	}{user, files}
	return c.JSON(http.StatusOK, response)
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

func (h fileHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file's id")
	}

	file, err := h.fileGetter.GetById(context.Background(), id.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	if err := h.fileDeleter.Delete(context.Background(), file); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
	}

	return c.NoContent(http.StatusOK)
}

func (h fileHandler) DeleteAllForOwner(c echo.Context) error {
	username := c.Param("username")

	user, err := h.userGetter.GetByUsername(context.Background(), username)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	if err := h.fileDeleter.DeleteAllForOwner(context.Background(), user); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
	}

	return c.NoContent(http.StatusOK)
}
