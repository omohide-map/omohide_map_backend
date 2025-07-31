package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/omohide_map_backend/internal/db"
	"github.com/omohide_map_backend/internal/handlers"
	omohideMiddleware "github.com/omohide_map_backend/internal/middleware"
	"github.com/omohide_map_backend/internal/models"
	"github.com/omohide_map_backend/internal/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	database := db.ConnectDB()
	log.Println("Database connected:", database != nil)

	if err := database.AutoMigrate(&models.Post{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Omohide Map API running!")
	})

	api := e.Group("/api")
	api.Use(omohideMiddleware.JWTMiddleware())

	postService := services.NewPostService(database)
	postHandler := handlers.NewPostHandler(postService)
	api.POST("/posts", postHandler.CreatePost)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}
