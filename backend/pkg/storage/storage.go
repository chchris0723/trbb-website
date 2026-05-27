package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"trbb/internal/config"
)

type Storage struct {
	client *minio.Client
	cfg    config.MinioConfig
}

func New(cfg config.MinioConfig) (*Storage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("init minio: %w", err)
	}

	s := &Storage{client: client, cfg: cfg}

	// Ensure default buckets exist
	for _, bucket := range []string{cfg.BucketPublic, cfg.BucketPrivate} {
		if err := s.ensureBucket(context.Background(), bucket); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Storage) ensureBucket(ctx context.Context, bucket string) error {
	exists, err := s.client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("check bucket %s: %w", bucket, err)
	}
	if !exists {
		if err := s.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("make bucket %s: %w", bucket, err)
		}
	}
	return nil
}

// Upload puts an object into the specified bucket and returns the public URL.
func (s *Storage) Upload(ctx context.Context, bucket, objectName, contentType string, reader io.Reader, size int64) (string, error) {
	_, err := s.client.PutObject(ctx, bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("upload object: %w", err)
	}

	url := fmt.Sprintf("%s/%s/%s", s.cfg.ExternalURL, bucket, objectName)
	return url, nil
}

// Delete removes an object from a bucket.
func (s *Storage) Delete(ctx context.Context, bucket, objectName string) error {
	return s.client.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
}

// PresignedURL generates a temporary presigned GET URL for private objects.
func (s *Storage) PresignedURL(ctx context.Context, bucket, objectName string, expiry time.Duration) (string, error) {
	u, err := s.client.PresignedGetObject(ctx, bucket, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("presign: %w", err)
	}
	return u.String(), nil
}

// PublicBucket returns the public bucket name.
func (s *Storage) PublicBucket() string { return s.cfg.BucketPublic }

// PrivateBucket returns the private bucket name.
func (s *Storage) PrivateBucket() string { return s.cfg.BucketPrivate }
