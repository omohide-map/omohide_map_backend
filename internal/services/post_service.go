// Package services contains business logic for the application.
package services

import (
	"time"

	"github.com/omohide_map_backend/internal/models"
	"github.com/omohide_map_backend/internal/utils"
	"gorm.io/gorm"
)

type PostService struct {
	db *gorm.DB
}

func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db: db}
}

func (s *PostService) CreatePost(req *models.CreatePostRequest, userID, userName, userAvatar string) (*models.Post, error) {
	post := &models.Post{
		ID:        utils.GenerateUlid(),
		UserID:    userID,
		Text:      req.Text,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		ImageUrls: []string{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(post).Error; err != nil {
		return nil, err
	}

	return post, nil
}
