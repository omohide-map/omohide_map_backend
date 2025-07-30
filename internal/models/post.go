// Package models contains the data models for the application.
package models

import (
	"time"
)

type Post struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"not null"`
	Text      string    `json:"text" gorm:"not null"`
	Latitude  float64   `json:"latitude" gorm:"not null"`
	Longitude float64   `json:"longitude" gorm:"not null"`
	ImageUrls []string  `json:"image_urls" gorm:"type:text[]"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

type CreatePostRequest struct {
	Text      string   `json:"text" validate:"required"`
	Latitude  float64  `json:"latitude" validate:"required"`
	Longitude float64  `json:"longitude" validate:"required"`
	Images    []string `json:"images"`
}
