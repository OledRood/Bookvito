package storage

import (
	"bookvito/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3ImageStorage struct {
	client        *minio.Client
	bucket        string
	publicBaseURL string
}

func NewS3ImageStorage(endpoint, accessKey, secretKey, bucket, publicBaseURL string, useSSL bool) (*S3ImageStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, err
		}
	}
	if err := setPublicReadPolicy(ctx, client, bucket); err != nil {
		return nil, err
	}

	return &S3ImageStorage{
		client:        client,
		bucket:        bucket,
		publicBaseURL: normalizePublicBaseURL(publicBaseURL),
	}, nil
}

func setPublicReadPolicy(ctx context.Context, client *minio.Client, bucket string) error {
	policy := map[string]any{
		"Version": "2012-10-17",
		"Statement": []map[string]any{
			{
				"Effect":    "Allow",
				"Principal": map[string]string{"AWS": "*"},
				"Action":    []string{"s3:GetObject"},
				"Resource":  []string{fmt.Sprintf("arn:aws:s3:::%s/*", bucket)},
			},
		},
	}
	raw, err := json.Marshal(policy)
	if err != nil {
		return err
	}
	return client.SetBucketPolicy(ctx, bucket, string(raw))
}

func (s *S3ImageStorage) Save(filename, contentType string, reader io.Reader) (string, error) {
	safeFilename, err := sanitizeFilename(filename)
	if err != nil {
		return "", err
	}

	_, err = s.client.PutObject(context.Background(), s.bucket, safeFilename, reader, -1, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	return s.PublicURL(safeFilename), nil
}

func (s *S3ImageStorage) Exists(publicURL string) (bool, error) {
	objectName, err := s.objectNameFromPublicURL(publicURL)
	if err != nil {
		return false, err
	}

	_, err = s.client.StatObject(context.Background(), s.bucket, objectName, minio.StatObjectOptions{})
	if err == nil {
		return true, nil
	}

	if minio.ToErrorResponse(err).Code == "NoSuchKey" {
		return false, nil
	}
	return false, err
}

func (s *S3ImageStorage) Delete(publicURL string) error {
	objectName, err := s.objectNameFromPublicURL(publicURL)
	if err != nil {
		return err
	}

	return s.client.RemoveObject(context.Background(), s.bucket, objectName, minio.RemoveObjectOptions{})
}

func (s *S3ImageStorage) PublicURL(filename string) string {
	return s.publicBaseURL + strings.TrimLeft(filename, "/")
}

func (s *S3ImageStorage) objectNameFromPublicURL(publicURL string) (string, error) {
	value := strings.TrimSpace(publicURL)
	if value == "" {
		return "", domain.NewValidationError("image_url пустой")
	}

	if strings.HasPrefix(value, s.publicBaseURL) {
		return sanitizeFilename(strings.TrimPrefix(value, s.publicBaseURL))
	}

	parsed, err := url.Parse(value)
	if err == nil && parsed.Path != "" {
		return sanitizeFilename(parsed.Path[strings.LastIndex(parsed.Path, "/")+1:])
	}

	return "", domain.NewValidationError("image_url не принадлежит image storage")
}

func NewImageStorage(driver, localDir, publicBaseURL, endpoint, accessKey, secretKey, bucket string, useSSL bool) (domain.ImageStorage, error) {
	switch strings.ToLower(strings.TrimSpace(driver)) {
	case "", "local":
		return NewLocalImageStorage(localDir, publicBaseURL)
	case "s3", "minio":
		return NewS3ImageStorage(endpoint, accessKey, secretKey, bucket, publicBaseURL, useSSL)
	default:
		return nil, fmt.Errorf("unsupported storage driver: %s", driver)
	}
}
