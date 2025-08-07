package di

import (
	"context"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/omohide_map_backend/internal/handler"
	"github.com/omohide_map_backend/internal/repository/repositories"
	"github.com/omohide_map_backend/internal/service"
	"github.com/omohide_map_backend/internal/storage"
	"google.golang.org/api/option"
)

type Container struct {
	AuthClient     *auth.Client
	FirestoreClient *firestore.Client
	S3Storage      *storage.S3Storage
	PostRepository *repositories.PostRepository
	PostService    *service.PostService
	PostHandler    *handler.PostHandler
}

func NewContainer(ctx context.Context) (*Container, error) {
	// Firebase初期化
	credentialsPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credentialsPath == "" {
		return nil, fmt.Errorf("GOOGLE_APPLICATION_CREDENTIALS environment variable is required")
	}

	opt := option.WithCredentialsFile(credentialsPath)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firebase app: %w", err)
	}

	// Auth Client
	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Auth client: %w", err)
	}

	// Firestore Client
	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Firestore client: %w", err)
	}

	// S3 Storage
	s3Storage, err := storage.NewS3Storage()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize S3 storage: %w", err)
	}

	// Repositories
	postRepo := repositories.NewPostRepository(firestoreClient)

	// Services
	postService := service.NewPostService(postRepo, s3Storage)

	// Handlers
	postHandler := handler.NewPostHandler(postService)

	return &Container{
		AuthClient:      authClient,
		FirestoreClient: firestoreClient,
		S3Storage:       s3Storage,
		PostRepository:  postRepo,
		PostService:     postService,
		PostHandler:     postHandler,
	}, nil
}

func (c *Container) Close() error {
	if c.FirestoreClient != nil {
		return c.FirestoreClient.Close()
	}
	return nil
}