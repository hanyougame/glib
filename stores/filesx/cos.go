package filesx

import (
	"context"
	"errors"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"net/http"
	"net/url"
)

// CosStorage 实现腾讯云
type CosStorage struct {
	client     *cos.Client
	requestURL string
}

type CosStorageConfig struct {
	// bucket地址 用于解析appid bucket_name
	BucketURL string `json:"bucket_url"`
	// 请求地址 用于请求图片
	RequestURL string `json:"request_url"`
	// app包请求地址
	PkgRequestUrl string `json:"pkg_request_url,optional"`
	// 用户 secret_id
	SecretID string `json:"secret_id"`
	// 密钥
	SecretKey string `json:"secret_key"`
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
