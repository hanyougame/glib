package filesx

import (
	"context"
	"errors"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/zeromicro/go-zero/core/stringx"
	"io"
	"net/http"
	"net/url"
	"time"
)

type CosSignType int // cos签名类型
const (
	_                       CosSignType = iota
	SignCosTypeAppLogUpload             // app日志上传
)

// CosStorage 实现腾讯云
type CosStorage struct {
	client     *cos.Client
	requestURL string
}

type CosStorageConfig struct {
	// bucket地址 用于解析appid bucket_name
	BucketURL string `json:"bucket_url"`
	// app包bucket
	PkgBucketURL string `json:"pkg_bucket_url,optional"`
	// 请求地址 用于请求图片
	RequestURL string `json:"request_url"`
	// app包请求地址
	PkgRequestUrl string `json:"pkg_request_url,optional"`
	// 用户 secret_id
	SecretID string `json:"secret_id"`
	// 密钥
	SecretKey string `json:"secret_key"`
	// 签名url过期时间，单位秒
	SignedURLExpire int64 `json:"signed_url_expire,optional"`
}

func NewCosStorage(c CosStorageConfig) (*CosStorage, error) {
	// 参数校验
	if c.SecretID == "" {
		return nil, errors.New("cos secret_id is require")
	}
	if c.SecretKey == "" {
		return nil, errors.New("cos secret_key is require")
	}
	if c.BucketURL == "" {
		return nil, errors.New("cos bucket_url is require")
	}

	u, err := url.Parse(c.BucketURL)
	if err != nil {
		return nil, fmt.Errorf("cos parse bucket_url fail:%s", err.Error())
	}
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  c.SecretID,  // 用户的 SecretId，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
			SecretKey: c.SecretKey, // 用户的 SecretKey，建议使用子账号密钥，授权遵循最小权限指引，降低使用风险。子账号密钥获取可参见 https://cloud.tencent.com/document/product/598/37140
		},
	})
	return &CosStorage{
		requestURL: c.RequestURL,
		client:     client,
	}, nil
}

func (s *CosStorage) Upload(ctx context.Context, file io.Reader, path string) (fullPath string, err error) {
	// 执行文件上传
	_, err = s.client.Object.Put(ctx, path, file, nil)
	if err != nil {
		err = fmt.Errorf("failed to upload file to cos: %w", err)
		return
	}
	fullPath, err = url.JoinPath(s.requestURL, path)
	return
}

func (s *CosStorage) Delete(ctx context.Context, path string) error {
	_, err := s.client.Object.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to delete file from cos: %w", err)
	}
	return nil
}

func (s *CosStorage) Exist(ctx context.Context, fileName string) (bool, error) {
	if _, err := s.client.Object.Head(ctx, fileName, nil); err != nil {
		if cos.IsNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GetPkgUploadSignUrl(c CosStorageConfig, key string) (string, error) {
	// 参数校验
	if stringx.HasEmpty(c.SecretID) {
		return "", fmt.Errorf("secret_id is require")
	}
	if stringx.HasEmpty(c.SecretKey) {
		return "", fmt.Errorf("secret_key is require")
	}
	if stringx.HasEmpty(c.PkgBucketURL) {
		return "", fmt.Errorf("bucket url is require")
	}
	if stringx.HasEmpty(key) {
		return "", fmt.Errorf("key is require")
	}
	u, err := url.Parse(c.PkgBucketURL)
	if err != nil {
		return "", fmt.Errorf("parse bucket url error: %+v", err)
	}

	// 初始化客户端
	client := cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  c.SecretID,
			SecretKey: c.SecretKey,
		},
	})

	// 获取预签名 URL
	presignedURL, err := client.Object.GetPresignedURL(context.Background(), http.MethodPut, key, c.SecretID, c.SecretKey, time.Hour, nil)
	if err != nil {
		return "", fmt.Errorf("GetPresignedURL error: %+v", err)
	}

	return presignedURL.String(), nil

}

func NewCosClient(secretID, secretKey, bucketURL string) (*cos.Client, error) {
	if secretID == "" {
		return nil, errors.New("cos secret_id is require")
	}
	if secretKey == "" {
		return nil, errors.New("cos secret_key is require")
	}
	if bucketURL == "" {
		return nil, errors.New("cos bucket_url is require")
	}

	u, err := url.Parse(bucketURL)
	if err != nil {
		return nil, fmt.Errorf("cos parse bucket_url fail:%s", err.Error())
	}

	return cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	}), nil
}

func (s CosSignType) ToString() string {
	switch s {
	case SignCosTypeAppLogUpload:
		return "app_log_upload"
	}
	return ""
}
