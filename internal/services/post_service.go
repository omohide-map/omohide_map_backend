// Package services contains business logic for the application.
package services

import (
	"github.com/omohide_map_backend/internal/models"
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
		Text:       req.Text,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		Images:     req.Images,
		ImageUrls:  []string{},
		UserID:     userID,
		UserName:   userName,
		UserAvatar: userAvatar,
	}

	if err := s.db.Create(post).Error; err != nil {
		return nil, err
	}

	return post, nil
}
