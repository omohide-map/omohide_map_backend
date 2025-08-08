package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/omohide_map_backend/internal/di"
)

func RegisterMainRoutes(api *echo.Group, container *di.Container) {
	api.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "JWT Token is valid")
	})

	api.GET("/posts", container.PostHandler.GetPosts)
	api.GET("/posts/my", container.PostHandler.GetMyPosts)
	// api.GET("/posts/:user_id", container.PostHandler.GetPostsByUserID)
	// api.GET("/post/:id", container.PostHandler.GetPostByID)
	api.POST("/post", container.PostHandler.CreatePost)
}
