package handler

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/omohide_map_backend/internal/models"
	"github.com/omohide_map_backend/internal/service"
	appErrors "github.com/omohide_map_backend/pkg/errors"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

func (h *PostHandler) CreatePost(c echo.Context) error {
	var req models.CreatePostRequest
	if err := c.Bind(&req); err != nil {
		return appErrors.InvalidRequest("Invalid request body")
	}
	if err := c.Validate(req); err != nil {
		return appErrors.ValidationError(err.Error())
	}

	userID, ok := c.Get("userID").(string)
	if !ok {
		return appErrors.UserIDNotFound()
	}

	requestTime, ok := c.Get("requestTime").(time.Time)
	if !ok {
		return appErrors.RequestTimeNotFound()
	}

	post, err := h.postService.CreatePost(c.Request().Context(), userID, &req, requestTime)
	if err != nil {
		return err
	}

	return c.JSON(201, post)
}
