package filesx

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/hanyougame/glib/stores/filesx/config"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	Local     = "local"
	AmazonS3  = "s3"
	AliyunOSS = "oss"
	Qiniu     = "qiniu"
)

// Storage 存储接口
type Storage interface {
	Upload(ctx context.Context, file io.Reader, path, contentType string) (string, error)
	Delete(ctx context.Context, path string) error
}

var (
	once     sync.Once
	Uploader *UploadManager
)

// UploadManager 上传管理器
type UploadManager struct {
	mu       sync.RWMutex
	configs  map[string]config.StorageConfig
	storages map[string]Storage
	errors   []error
}

// Must 初始化上传管理器
func Must(configs ...config.StorageConfig) {
	var err error

	once.Do(func() {
		Uploader = &UploadManager{
			configs:  make(map[string]config.StorageConfig),
			storages: make(map[string]Storage),
			errors:   make([]error, 0),
		}
		for _, c := range configs {
			if err = Uploader.addStorage(c); err != nil {
				Uploader.errors = append(Uploader.errors, err)
			}
		}
	})

	if len(Uploader.errors) > 0 {
		panic(fmt.Sprintf("errors initializing storages: %v", Uploader.errors))
	}
}

func (u *UploadManager) addStorage(cfg config.StorageConfig) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if cfg.Type == "" || cfg.Bucket == "" {
		return fmt.Errorf("invalid storage configuration for type %s", cfg.Type)
	}

	u.configs[cfg.Type] = cfg
	switch cfg.Type {
	case AmazonS3:
		s, err := newS3Storage(cfg)
		if err != nil {
			return err
		}
		u.storages[cfg.Type] = s
	case Local:
		s, err := newLocalStorage(cfg)
		if err != nil {
			return err
		}
		u.storages[cfg.Type] = s
	default:
		return fmt.Errorf("unsupported storage type: %s", cfg.Type)
	}
	return nil
}

// Upload 上传文件
func (u *UploadManager) Upload(ctx context.Context, storageType string, file io.Reader, header *multipart.FileHeader) (string, error) {
	storage, ok := u.storages[storageType]
	if !ok {
		return "", fmt.Errorf("storage type %s not initialized", storageType)
	}

	// 根据日期生成文件路径
	datePath := time.Now().Format("2006/01/02")
	fileExt := filepath.Ext(header.Filename)
	uniqueFilename := fmt.Sprintf("%s%s", uuid.New().String(), fileExt)

	// 检测 Content-Type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		var err error
		contentType, err = u.detectContentType(file, header.Filename)
		if err != nil {
			return "", fmt.Errorf("failed to detect content type: %w", err)
		}
	}

	// 重置文件读取位置（如果文件是 io.Seeker）
	if seeker, ok := file.(io.Seeker); ok {
		_, err := seeker.Seek(0, io.SeekStart)
		if err != nil {
			return "", fmt.Errorf("failed to reset file reader: %w", err)
		}
	}

	// 生成文件上传路径
	path := filepath.Join(datePath, uniqueFilename)

	// 执行上传操作
	url, err := storage.Upload(ctx, file, path, contentType)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return url, nil
}

// Delete 删除文件
func (u *UploadManager) Delete(ctx context.Context, storageType string, path string) error {
	storage, ok := u.storages[storageType]
	if !ok {
		return fmt.Errorf("storage type %s not initialized", storageType)
	}
	return storage.Delete(ctx, path)
}

func (u *UploadManager) detectContentType(file io.Reader, filename string) (string, error) {
	// 尝试通过文件扩展名来确定 MIME 类型
	ext := strings.ToLower(filepath.Ext(filename))
	ct := mime.TypeByExtension(ext)
	if ct != "" {
		return ct, nil
	}

	// 读取文件的前 512 字节来检测 MIME 类型
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("failed to read file for content type detection: %w", err)
	}

	// 如果成功读取了数据，检测 MIME 类型
	if n > 0 {
		return http.DetectContentType(buf[:n]), nil
	}

	// 如果读取的字节数为 0，无法检测类型，返回默认值
	return "application/octet-stream", nil
}
