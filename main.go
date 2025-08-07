package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/omohide_map_backend/internal/di"
	omohideMiddleware "github.com/omohide_map_backend/internal/middleware"
	"github.com/omohide_map_backend/pkg/validator"
)

func main() {
	ctx := context.Background()

	// .envファイル読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// DIコンテナの初期化
	container, err := di.NewContainer(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize DI container: %v", err)
	}

	// echo
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = validator.New()
	e.HTTPErrorHandler = omohideMiddleware.CustomErrorHandler

	// sample
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Omohide Map API running!")
	})

	api := e.Group("/api")
	api.Use(omohideMiddleware.JWTMiddleware(container.AuthClient))

	api.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "JWT Token is valid")
	})

	// endpoints
	api.POST("/post", container.PostHandler.CreatePost)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// グレースフルシャットダウンの設定
	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// シグナル待機
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	// グレースフルシャットダウン
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

	// DIコンテナのクリーンアップ
	if err := container.Close(); err != nil {
		log.Printf("Failed to close DI container: %v", err)
	}

	log.Println("Server gracefully stopped")
}
