package models

import (
	"time"

	"github.com/lib/pq"
)

type Post struct {
	ID        string         `json:"id" firestore:"id" gorm:"primaryKey"`
	UserID    string         `json:"user_id" firestore:"user_id" gorm:"not null"`
	Text      string         `json:"text" firestore:"text" gorm:"not null"`
	Latitude  float64        `json:"latitude" firestore:"latitude" gorm:"not null"`
	Longitude float64        `json:"longitude" firestore:"longitude" gorm:"not null"`
	ImageUrls pq.StringArray `json:"image_urls" firestore:"image_urls" gorm:"type:text[]"`
	CreatedAt time.Time      `json:"created_at" firestore:"created_at" gorm:"not null"`
	UpdatedAt time.Time      `json:"updated_at" firestore:"updated_at" gorm:"not null"`
}

type CreatePostRequest struct {
	Text      string   `json:"text" validate:"required"`
	Latitude  float64  `json:"latitude" validate:"required"`
	Longitude float64  `json:"longitude" validate:"required"`
	Images    []string `json:"images"`
}

type GetPostsRequest struct {
	Page      *int     `query:"page"`
	Limit     *int     `query:"limit"`
	Latitude  *float64 `query:"latitude"`
	Longitude *float64 `query:"longitude"`
	Radius    *float64 `query:"radius"`
}
