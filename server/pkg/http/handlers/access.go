package handlers

import (
	"context"
	"net/http"

	"github.com/SergeyCherepiuk/docs/pkg/database/models"
	"github.com/SergeyCherepiuk/docs/pkg/database/neo4j"
	"github.com/SergeyCherepiuk/docs/pkg/http/internal"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AccessHandler struct{}

func (h AccessHandler) Grant(c echo.Context) error {
	user, ok := c.Get("user").(models.User)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "User wasn't found")
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file id")
	}

	ctx := context.Background()
	sess := neo4j.NewSession(ctx)
	defer sess.Close(ctx)

	file, err := neo4j.FileService.GetById(ctx, sess, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, internal.ToSentence(err.Error()))
	}

	type RequestBody struct {
		Receiver string `json:"receiver"`
		Level    string `json:"level"`
	}

	var body RequestBody
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	// TODO: Validation

	receiver, err := neo4j.UserService.GetByUsername(ctx, sess, body.Receiver)
	if err != nil {
		return err
	}

	access := models.Access{
		Granter:  user.Username,
		Receiver: receiver.Username,
		Level:    body.Level,
	}

	prevAccess, prevAccessErr := neo4j.AccessService.Get(ctx, sess, file, receiver)
	if prevAccessErr != nil {
		err = neo4j.AccessService.Grant(ctx, sess, file, access)
	} else {
		err = neo4j.AccessService.UpdateLevel(ctx, sess, file, prevAccess, access.Level)
	}

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, internal.ToSentence(err.Error()))
	}

	return c.NoContent(http.StatusCreated)
}

func (h AccessHandler) GetAccesses(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid file id")
	}

	ctx := context.Background()
	sess := neo4j.NewSession(ctx)
	defer sess.Close(ctx)

	file, err := neo4j.FileService.GetById(ctx, sess, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	accesses, err := neo4j.AccessService.GetAccesses(ctx, sess, file)
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

	ctx := context.Background()
	sess := neo4j.NewSession(ctx)
	defer sess.Close(ctx)

	file, err := neo4j.FileService.GetById(ctx, sess, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	username := c.Param("username")
	user, err := neo4j.UserService.GetByUsername(ctx, sess, username)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	access, err := neo4j.AccessService.Get(ctx, sess, file, user)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	if err := neo4j.AccessService.Revoke(ctx, sess, file, access); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, internal.ToSentence(err.Error()))
	}

	return c.NoContent(http.StatusOK)
}
