package repository

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/omohide_map_backend/internal/models"
	appErrors "github.com/omohide_map_backend/pkg/errors"
	"github.com/omohide_map_backend/pkg/geo"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		if status.Code(err) == codes.NotFound {
			return nil, appErrors.ResourceNotFound("Post")
		}
		return nil, appErrors.DatabaseError(err)
	}

	var post models.Post
	if err := doc.DataTo(&post); err != nil {
		return nil, appErrors.DatabaseError(err)
	}

	return &post, nil
}

func (r *PostRepository) GetByUserID(ctx context.Context, userID string) ([]*models.Post, error) {
	iter := r.firestoreClient.Collection("posts").Where("user_id", "==", userID).Documents(ctx)
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

func (r *PostRepository) GetPostsWithFilters(ctx context.Context, req *models.GetPostsRequest) ([]*models.Post, error) {
	log.Printf("GetPostsWithFilters called with req: %+v", req)

	// OrderByを復活させる
	query := r.firestoreClient.Collection("posts").OrderBy("created_at", firestore.Desc)

	// ページネーション
	limit := 20 // デフォルト
	if req.Limit != nil && *req.Limit > 0 && *req.Limit <= 100 {
		limit = *req.Limit
	}

	offset := 0
	if req.Page != nil && *req.Page > 0 {
		offset = (*req.Page - 1) * limit
	}

	log.Printf("Pagination: limit=%d, offset=%d", limit, offset)
	query = query.Limit(limit).Offset(offset)

	iter := query.Documents(ctx)
	defer iter.Stop()

	var posts []*models.Post
	docCount := 0
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			log.Printf("Iterator done. Total documents fetched: %d", docCount)
			break
		}
		if err != nil {
			log.Printf("Error fetching document: %v", err)
			return nil, err
		}
		docCount++

		var post models.Post
		if err := doc.DataTo(&post); err != nil {
			log.Printf("Error converting document to Post struct: %v", err)
			continue
		}

		log.Printf("Fetched post: ID=%s, UserID=%s", post.ID, post.UserID)

		// 位置フィルタリング（メモリ内で処理）
		// 緯度・経度・半径が全て有効な値の場合のみフィルタリング
		if req.Latitude != nil && req.Longitude != nil && req.Radius != nil && *req.Radius > 0 {
			distance := geo.CalculateDistance(*req.Latitude, *req.Longitude, post.Latitude, post.Longitude)
			if distance > *req.Radius {
				log.Printf("Post %s filtered out by distance: %f > %f", post.ID, distance, *req.Radius)
				continue
			}
		}

		posts = append(posts, &post)
	}

	log.Printf("Returning %d posts", len(posts))
	return posts, nil
}
