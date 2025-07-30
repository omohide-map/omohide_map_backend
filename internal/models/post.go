// Package models contains the data models for the application.
package models

import (
	"time"
)

type Post struct {
	ID         string    `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Text       string    `json:"text" gorm:"not null"`
	Latitude   float64   `json:"latitude" gorm:"not null"`
	Longitude  float64   `json:"longitude" gorm:"not null"`
	Images     []string  `json:"images" gorm:"type:text[]"`
	ImageUrls  []string  `json:"imageUrls" gorm:"type:text[]"`
	CreatedAt  time.Time `json:"createdAt"`
	UserID     string    `json:"userId" gorm:"not null"`
	UserName   string    `json:"userName" gorm:"not null"`
	UserAvatar string    `json:"userAvatar"`
}

type CreatePostRequest struct {
	Text      string    `json:"text" validate:"required"`
	Latitude  float64   `json:"latitude" validate:"required"`
	Longitude float64   `json:"longitude" validate:"required"`
	Images    []string  `json:"images"`
}
