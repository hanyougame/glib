package filex

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Config 上传配置
type Config struct {
	// 存储配置
	Storage StorageConfig `json:"storage"`
	// 上传限制配置
	Upload UploadConfig `json:"upload"`
}

// StorageConfig 存储配置
type StorageConfig struct {
	// 存储类型: local, s3, oss 等
	Type string `json:"type"`
	// 本地存储路径或云存储 bucket
	Bucket string `json:"bucket"`
	// 访问密钥
	AccessKey string `json:"access_key"`
	// 密钥
	SecretKey string `json:"secret_key"`
	// 区域
	Region string `json:"region"`
	// CDN域名（可选）
	CdnDomain string `json:"cdn_domain"`
}

// UploadConfig 上传限制配置
type UploadConfig struct {
	// 最大文件大小（字节）
	MaxSize int64 `json:"max_size"`
	// 允许的文件类型
	AllowedTypes []string `json:"allowed_types"`
	// 基础上传路径
	BasePath string `json:"base_path"`
}

// FileInfo 上传文件信息
type FileInfo struct {
	URL      string `json:"url"`       // 完整访问URL
	Path     string `json:"path"`      // 相对存储路径
	Name     string `json:"name"`      // 文件名
	Size     int64  `json:"size"`      // 文件大小
	MimeType string `json:"mime_type"` // 文件类型
}

// Storage 存储接口
type Storage interface {
	Upload(ctx context.Context, file io.Reader, path string) (string, error)
	Delete(ctx context.Context, path string) error
}

// Uploader 上传管理器
type Uploader struct {
	config  Config
	storage Storage
	once    sync.Once
}

// New 创建上传管理器
func New(config Config) (*Uploader, error) {
	u := &Uploader{
		config: config,
	}

	// 根据配置初始化存储
	var err error
	u.once.Do(func() {
		switch config.Storage.Type {
		case "s3":
			u.storage, err = newS3Storage(config.Storage)
		case "local":
			u.storage, err = newLocalStorage(config.Storage)
		default:
			err = fmt.Errorf("unsupported storage type: %s", config.Storage.Type)
		}
	})

	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	return u, nil
}

// Upload 上传文件
func (u *Uploader) Upload(ctx context.Context, file *multipart.FileHeader) (*FileInfo, error) {
	// 验证文件
	if err := u.validateFile(file); err != nil {
		return nil, err
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// 生成存储路径
	filename := u.generateFileName(file.Filename)
	path := u.generatePath(filename)

	// 上传文件
	url, err := u.storage.Upload(ctx, src, path)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	return &FileInfo{
		URL:      url,
		Path:     path,
		Name:     filename,
		Size:     file.Size,
		MimeType: file.Header.Get("Content-Type"),
	}, nil
}

// validateFile 验证文件
func (u *Uploader) validateFile(file *multipart.FileHeader) error {
	if file.Size > u.config.Upload.MaxSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %d bytes", u.config.Upload.MaxSize)
	}

	mimeType := file.Header.Get("Content-Type")
	if !u.isAllowedType(mimeType) {
		return fmt.Errorf("file type %s is not allowed", mimeType)
	}

	return nil
}

// isAllowedType 检查文件类型是否允许
func (u *Uploader) isAllowedType(mimeType string) bool {
	for _, allowed := range u.config.Upload.AllowedTypes {
		if strings.HasPrefix(mimeType, allowed) {
			return true
		}
	}
	return false
}

// generateFileName 生成文件名
func (u *Uploader) generateFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	return fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
}

// generatePath 生成存储路径
func (u *Uploader) generatePath(filename string) string {
	return filepath.Join(
		u.config.Upload.BasePath,
		time.Now().Format("2006/01/02"),
		filename,
	)
}
