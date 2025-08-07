package main

import (
	"context"
	"log"
	"net/http"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	omohideMiddleware "github.com/omohide_map_backend/internal/middleware"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	// .envファイル読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	// firebase
	opt := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}
	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("app.Auth: %v", err)
	}
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Omohide Map API running!")
	})

	api := e.Group("/api")
	api.Use(omohideMiddleware.JWTMiddleware(authClient))

	api.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "JWT Token is valid")
	})

	api.POST("/posts", postHandler.CreatePost) // TODO: あとで実装

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
	defer firestoreClient.Close()
}
