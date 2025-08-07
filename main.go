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
	"github.com/omohide_map_backend/internal/handler"
	omohideMiddleware "github.com/omohide_map_backend/internal/middleware"
	"github.com/omohide_map_backend/internal/repository"
	"github.com/omohide_map_backend/internal/service"
	"github.com/omohide_map_backend/internal/storage"
	"github.com/omohide_map_backend/pkg/validator"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()

	// .envファイル読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	// firebase
	/// Authentication
	opt := option.WithCredentialsFile(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v", err)
	}
	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("app.Auth: %v", err)
	}

	/// Firestore
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("app.Firestore: %v", err)
	}

	// S3 Storage
	bucketName := os.Getenv("AWS_S3_BUCKET")
	if bucketName == "" {
		log.Fatal("AWS_S3_BUCKET environment variable is required")
	}
	s3Storage, err := storage.NewS3Storage(bucketName)
	if err != nil {
		log.Fatalf("Failed to initialize S3 storage: %v", err)
	}

	// handler
	postRepo := repository.NewPostRepository(firestoreClient)
	postService := service.NewPostService(postRepo, s3Storage)
	postHandler := handler.NewPostHandler(postService)

	// echo
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Validator = validator.New()

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Omohide Map API running!")
	})

	api := e.Group("/api")
	api.Use(omohideMiddleware.JWTMiddleware(authClient))

	api.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "JWT Token is valid")
	})

	api.POST("/post", postHandler.CreatePost)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
	defer firestoreClient.Close()
}
