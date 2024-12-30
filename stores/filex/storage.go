package filex

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"os"
	"path/filepath"
)

// s3Storage 实现AWS S3的存储接口
type s3Storage struct {
	client    *s3.Client
	bucket    string
	region    string
	cdnDomain string // CDN域名（可选）
}

func newS3Storage(sc StorageConfig) (Storage, error) {
	if sc.AccessKey == "" || sc.SecretKey == "" {
		return nil, fmt.Errorf("s3 credentials are required")
	}
	if sc.Bucket == "" {
		return nil, fmt.Errorf("s3 bucket is required")
	}
	if sc.Region == "" {
		return nil, fmt.Errorf("s3 region is required")
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(sc.Region))
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS SDK config: %w", err)
	}

	if sc.AccessKey != "" && sc.SecretKey != "" {
		cfg.Credentials = credentials.NewStaticCredentialsProvider(sc.AccessKey, sc.SecretKey, "")
	}

	client := s3.NewFromConfig(cfg)
	return &s3Storage{
		client:    client,
		bucket:    sc.Bucket,
		region:    sc.Region, // Ensure region is set
		cdnDomain: sc.CdnDomain,
	}, nil
}
func (s *s3Storage) Upload(ctx context.Context, file io.Reader, path string) (string, error) {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
		Body:   file,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Construct the URL based on whether a CDN domain is provided
	var url string
	if s.cdnDomain != "" {
		url = fmt.Sprintf("https://%s/%s", s.cdnDomain, path)
	} else {
		url = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, path)
	}

	return url, nil
}

func (s *s3Storage) Delete(ctx context.Context, path string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}
	return nil
}

type localStorage struct {
	basePath  string
	cdnDomain string // CDN域名（可选）
}

func newLocalStorage(config StorageConfig) (Storage, error) {
	if config.Bucket == "" {
		return nil, fmt.Errorf("local storage path is required")
	}

	// Ensure the storage directory exists
	if err := os.MkdirAll(config.Bucket, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &localStorage{
		basePath:  config.Bucket,
		cdnDomain: config.CdnDomain,
	}, nil
}

func (l *localStorage) Upload(ctx context.Context, file io.Reader, path string) (string, error) {
	// Create full path including base path
	fullPath := filepath.Join(l.basePath, path)
	dir := filepath.Dir(fullPath)

	// Ensure the directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// Write the file to disk
	f, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, file); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// Construct the URL based on whether a CDN domain is provided
	var url string
	if l.cdnDomain != "" {
		url = fmt.Sprintf("https://%s/%s", l.cdnDomain, path)
	} else {
		url = fmt.Sprintf("file://%s", fullPath)
	}

	return url, nil
}

// Delete deletes a file from the local filesystem.
func (l *localStorage) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(l.basePath, path)
	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
