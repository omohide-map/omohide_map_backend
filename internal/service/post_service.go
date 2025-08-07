package service

import (
	"context"
	"fmt"
	"time"

	"github.com/omohide_map_backend/internal/models"
	"github.com/omohide_map_backend/internal/repository/repositories"
	"github.com/omohide_map_backend/internal/storage"
	"github.com/omohide_map_backend/pkg/id"
)

type PostService struct {
	postRepo  *repositories.PostRepository
	s3Storage *storage.S3Storage
}

func NewPostService(postRepo *repositories.PostRepository, s3Storage *storage.S3Storage) *PostService {
	return &PostService{
		postRepo:  postRepo,
		s3Storage: s3Storage,
	}
}

func (s *PostService) CreatePost(ctx context.Context, userID string, req *models.CreatePostRequest, requestTime time.Time) (*models.Post, error) {
	postID := id.GenerateUlid(requestTime)

	var imageURLs []string
	if len(req.Images) > 0 {
		for i, imageBase64 := range req.Images {
			key := fmt.Sprintf("posts/%s/%d.jpg", postID, i)
			url, err := s.s3Storage.UploadBase64Image(ctx, key, imageBase64)
			if err != nil {
				return nil, err
			}
			imageURLs = append(imageURLs, url)
		}
	}

	post := &models.Post{
		ID:        postID,
		UserID:    userID,
		Text:      req.Text,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		ImageUrls: imageURLs,
		CreatedAt: requestTime,
		UpdatedAt: requestTime,
	}

	if err := s.postRepo.Create(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}
