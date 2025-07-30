// Package handlers contains HTTP handlers for the application.
package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/omohide_map_backend/internal/models"
	"github.com/omohide_map_backend/internal/services"
)

type PostHandler struct {
	postService *services.PostService
}

func NewPostHandler(postService *services.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

func (h *PostHandler) CreatePost(c echo.Context) error {
	var req models.CreatePostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "User ID not found in token",
		})
	}

	userEmail, _ := c.Get("user_email").(string)
	userName := userEmail
	if userName == "" {
		userName = "Unknown User"
	}

	userAvatar := ""

	post, err := h.postService.CreatePost(&req, userID, userName, userAvatar)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create post",
		})
	}

	return c.JSON(http.StatusCreated, post)
}