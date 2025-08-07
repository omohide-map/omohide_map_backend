package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/omohide_map_backend/internal/models"
	"github.com/omohide_map_backend/internal/service"
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
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userID, ok := c.Get("userID").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "User ID not found")
	}

	requestTime, ok := c.Get("requestTime").(time.Time)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, "Request time not found")
	}

	post, err := h.postService.CreatePost(c.Request().Context(), userID, &req, requestTime)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, post)
}
