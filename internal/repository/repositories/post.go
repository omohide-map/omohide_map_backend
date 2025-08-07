package repositories

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/omohide_map_backend/internal/models"
)

type PostRepository struct {
	firestoreClient *firestore.Client
}

func NewPostRepository(firestoreClient *firestore.Client) *PostRepository {
	return &PostRepository{
		firestoreClient: firestoreClient,
	}
}

func (r *PostRepository) Create(ctx context.Context, post *models.Post) error {
	_, err := r.firestoreClient.Collection("posts").Doc(post.ID).Set(ctx, post)
	return err
}

func (r *PostRepository) GetByID(ctx context.Context, id string) (*models.Post, error) {
	doc, err := r.firestoreClient.Collection("posts").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	var post models.Post
	if err := doc.DataTo(&post); err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) GetByUserID(ctx context.Context, userID string) ([]*models.Post, error) {
	iter := r.firestoreClient.Collection("posts").Where("userID", "==", userID).Documents(ctx)
	defer iter.Stop()

	var posts []*models.Post
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}

		var post models.Post
		if err := doc.DataTo(&post); err != nil {
			continue
		}
		posts = append(posts, &post)
	}

	return posts, nil
}
