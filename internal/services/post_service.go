// Package services contains business logic for the application.
package services

import (
	"time"

	"github.com/lib/pq"
	"github.com/omohide_map_backend/internal/models"
	"github.com/omohide_map_backend/internal/storage"
	"github.com/omohide_map_backend/internal/utils"
	"gorm.io/gorm"
)

type PostService struct {
	db      *gorm.DB
	storage *storage.SupabaseStorage
}

func NewPostService(db *gorm.DB) *PostService {
	storage, err := storage.NewSupabaseStorage()
	if err != nil {
		panic("Failed to initialize Supabase storage: " + err.Error())
	}
	return &PostService{db: db, storage: storage}
}

func (s *PostService) CreatePost(req *models.CreatePostRequest, userID, userName, userAvatar string) (*models.Post, error) {
	imageUrls := []string{}
	bucketName := "posts"

	for _, base64Image := range req.Images {
		if base64Image != "" {
			imageURL, err := s.storage.UploadImage(base64Image, bucketName, userID)
			if err != nil {
				return nil, err
			}
			imageUrls = append(imageUrls, imageURL)
		}
	}

	post := &models.Post{
		ID:        utils.GenerateUlid(),
		UserID:    userID,
		Text:      req.Text,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		ImageUrls: pq.StringArray(imageUrls),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(post).Error; err != nil {
		return nil, err
	}

	return post, nil
}
