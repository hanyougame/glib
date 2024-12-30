package config

// StorageConfig 存储配置
type StorageConfig struct {
	// 存储类型: local, s3, oss 等
	Type string `json:"type,optional"`
	// 本地存储路径或云存储 bucket
	Bucket string `json:"bucket,optional"`
	// 访问密钥
	AccessKey string `json:"access_key,optional"`
	// 密钥
	SecretKey string `json:"secret_key,optional"`
	// 区域
	Region string `json:"region,optional"`
	// CDN域名（可选）
	CdnDomain string `json:"cdn_domain,optional"`
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
