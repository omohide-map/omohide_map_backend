// Package storage handles file storage operations.
package storage

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/omohide_map_backend/internal/utils"
)

type SupabaseStorage struct {
	storageURL string
	serviceKey string
}

func NewSupabaseStorage() (*SupabaseStorage, error) {
	supabaseURL := os.Getenv("SUPABASE_STORAGE_URL")
	supabaseServiceKey := os.Getenv("SUPABASE_API_SECRET_KEY")

	if supabaseURL == "" || supabaseServiceKey == "" {
		return nil, fmt.Errorf("SUPABASE_URL or SUPABASE_API_SECRET_KEY is not set")
	}

	return &SupabaseStorage{
		storageURL: supabaseURL,
		serviceKey: supabaseServiceKey,
	}, nil
}

func (s *SupabaseStorage) UploadImage(base64Data string, bucketName string, userID string) (string, error) {
	if strings.Contains(base64Data, ",") {
		base64Data = strings.Split(base64Data, ",")[1]
	}

	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 data: %w", err)
	}

	fileName := fmt.Sprintf("%s/%s.jpg", userID, utils.GenerateUlid())
	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", s.storageURL, bucketName, fileName)

	req, err := http.NewRequest("PUT", uploadURL, bytes.NewReader(imageData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.serviceKey))
	req.Header.Set("Content-Type", "image/jpeg")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(imageData)))
	req.Header.Set("apikey", s.serviceKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.storageURL, bucketName, fileName)

	return publicURL, nil
}
