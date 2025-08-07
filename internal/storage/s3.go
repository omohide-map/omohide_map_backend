package storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	appErrors "github.com/omohide_map_backend/pkg/errors"
)

type S3Storage struct {
	client     *s3.Client
	bucketName string
}

func NewS3Storage() (*S3Storage, error) {
	bucketName := os.Getenv("AWS_S3_BUCKET")
	if bucketName == "" {
		return nil, appErrors.EnvironmentVariableError("AWS_S3_BUCKET")
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		return nil, appErrors.StorageError(err)
	}

	client := s3.NewFromConfig(cfg)
	return &S3Storage{
		client:     client,
		bucketName: bucketName,
	}, nil
}

func (s *S3Storage) UploadBase64Image(ctx context.Context, key string, base64Image string) (string, error) {
	base64Data := base64Image
	if strings.Contains(base64Image, ",") {
		parts := strings.Split(base64Image, ",")
		if len(parts) == 2 {
			base64Data = parts[1]
		}
	}

	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", appErrors.ImageProcessingError(err)
	}

	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		return "", appErrors.StorageError(err)
	}

	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s.bucketName, key)
	return url, nil
}
