package core

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/hanyougame/glib/flowcore/config"
	"github.com/hanyougame/glib/flowcore/logger"
	"go.temporal.io/sdk/log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
)

var (
	defaultClient client.Client
	clientOnce    sync.Once
	clientMu      sync.RWMutex
	namedClients  = make(map[string]client.Client)
)

// todo 后续随着业务的增加单一的命名空间不太友好，考虑使用命名空间生成不同的客户端 并使用链接复用利用资源NewClientFromExistingWithContext
// GetOrCreateDefault 返回默认客户端，如果不存在则创建一个新的
func GetOrCreateDefault(ctx context.Context, cfg *config.Config) (client.Client, error) {
	var err error
	clientOnce.Do(func() {
		defaultClient, err = NewClient(ctx, cfg)
	})
	return defaultClient, err
}

// NewClient 创建一个新的 Temporal 客户端
func NewClient(ctx context.Context, cfg *config.Config) (client.Client, error) {
	if !cfg.Enabled {
		return nil, fmt.Errorf("temporal client is disabled")
	}

	// 🔧 如果启用自动创建命名空间，使用独立的命名空间管理器
	if cfg.NamespaceConf.AutoCreate {
		err := EnsureNamespaceExists(ctx, cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to ensure namespace exists: %w", err)
		}
	}

	options := buildClientOptions(cfg)

	// 配置日志
	if options.Logger == nil {
		if os.Getenv("ENV") == "production" {
			options.Logger = logger.NewProductionLogger()
		} else {
			options.Logger = logger.NewDevelopmentLogger("info")
		}
	}

	return client.DialContext(ctx, options)
}

// NewClientWithLogger 使用自定义日志记录器创建一个新的客户端
func NewClientWithLogger(ctx context.Context, cfg *config.Config, logger log.Logger) (client.Client, error) {
	if !cfg.Enabled {
		return nil, fmt.Errorf("temporal client is disabled")
	}

	options := buildClientOptions(cfg)
	options.Logger = logger

	return client.DialContext(ctx, options)
}

// buildClientOptions 构建客户端选项
func buildClientOptions(cfg *config.Config) client.Options {
	options := client.Options{
		HostPort:  cfg.HostPort,
		Namespace: cfg.Namespace,
	}

	// 设置身份标识
	if cfg.Identity != "" {
		options.Identity = cfg.Identity
	} else {
		hostname, _ := os.Hostname()
		options.Identity = fmt.Sprintf("temporal-worker@%s", hostname)
	}

	// 配置 TLS
	if cfg.TLS.Enabled {
		if tlsConfig, err := cfg.TLS.GetTLSConfig(); err == nil {
			options.ConnectionOptions = client.ConnectionOptions{
				TLS: tlsConfig,
			}
		}
	}

	// 配置数据转换器
	switch cfg.DataConverter {
	case "proto":
		options.DataConverter = converter.NewCompositeDataConverter(
			converter.NewProtoPayloadConverter(),
			converter.NewJSONPayloadConverter(),
		)
	default:
		options.DataConverter = converter.GetDefaultDataConverter()
	}

	return options
}

// RegisterClient 注册一个带有名称的客户端
func RegisterClient(name string, c client.Client) {
	clientMu.Lock()
	defer clientMu.Unlock()
	namedClients[name] = c
}

// GetClient 获取指定名称的客户端
func GetClient(name string) (client.Client, bool) {
	clientMu.RLock()
	defer clientMu.RUnlock()
	c, ok := namedClients[name]
	return c, ok
}

// MustGetClient 获取指定名称的客户端，如果不存在则触发 panic
func MustGetClient(name string) client.Client {
	if c, ok := GetClient(name); ok {
		return c
	}
	panic(fmt.Sprintf("client '%s' not found", name))
}

// CloseAll 关闭所有客户端
func CloseAll() {
	clientMu.Lock()
	defer clientMu.Unlock()

	if defaultClient != nil {
		defaultClient.Close()
		defaultClient = nil
	}

	for name, c := range namedClients {
		c.Close()
		delete(namedClients, name)
	}
}
