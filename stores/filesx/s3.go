package filesx

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	config2 "github.com/hanyougame/glib/stores/filesx/config"
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

func newS3Storage(sc config2.StorageConfig) (Storage, error) {
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

func (s *s3Storage) Upload(ctx context.Context, file io.Reader, path, contentType string) (string, error) {
	// 执行文件上传
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(path),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	if s.cdnDomain != "" {
		return fmt.Sprintf("https://%s/%s", s.cdnDomain, path), nil
	}
	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucket, s.region, path), nil
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

func (s *s3Storage) Exist(ctx context.Context, fileName string) (bool, error) {
	if _, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fileName),
	}); err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "NotFound" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

type localStorage struct {
	basePath  string
	cdnDomain string // CDN域名（可选）
}

func newLocalStorage(config config2.StorageConfig) (Storage, error) {
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

func (l *localStorage) Upload(ctx context.Context, file io.Reader, path, contentType string) (string, error) {
	// 创建完整的路径，包括基础路径和日期路径
	fullPath := filepath.Join(l.basePath, path)
	dir := filepath.Dir(fullPath)

	// 确保目录存在
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// 写入文件到磁盘
	f, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, file); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

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

func (l *localStorage) Exist(ctx context.Context, fileName string) (bool, error) {
	_, err := os.Stat(filepath.Join(l.basePath, fileName))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
