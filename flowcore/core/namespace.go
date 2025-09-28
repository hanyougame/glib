package core

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hanyougame/glib/flowcore/config"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"google.golang.org/protobuf/types/known/durationpb"
)

const (
	defaultRetentionPeriod = 72 * time.Hour // 默认保留期3天
)

// EnsureNamespaceExists 确保命名空间存在
func EnsureNamespaceExists(ctx context.Context, cfg *config.Config) error {
	nsClient, err := createNamespaceClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create namespace client: %w", err)
	}
	defer nsClient.Close()

	// 检查命名空间是否存在
	if _, err = nsClient.Describe(ctx, cfg.Namespace); err == nil {
		return nil // 命名空间已存在
	} else if !isNamespaceNotFoundError(err) {
		return err // 其他错误
	}

	// 创建命名空间
	return createNamespace(ctx, nsClient, cfg)
}

// createNamespace 创建命名空间
func createNamespace(ctx context.Context, nsClient client.NamespaceClient, cfg *config.Config) error {
	retentionPeriod := parseRetentionPeriod(cfg.NamespaceConf.RetentionPeriod)

	description := cfg.NamespaceConf.Description
	if description == "" {
		description = fmt.Sprintf("Namespace for %s", cfg.Namespace)
	}

	request := &workflowservice.RegisterNamespaceRequest{
		Namespace:                        cfg.Namespace,
		Description:                      description,
		WorkflowExecutionRetentionPeriod: retentionPeriod,
		IsGlobalNamespace:                false,
		HistoryArchivalState:             getArchivalState(cfg.NamespaceConf.HistoryArchivalEnabled),
		VisibilityArchivalState:          getArchivalState(cfg.NamespaceConf.VisibilityArchivalEnabled),
	}

	return nsClient.Register(ctx, request)
}

// createNamespaceClient 创建命名空间客户端
func createNamespaceClient(cfg *config.Config) (client.NamespaceClient, error) {
	options := client.Options{
		HostPort: cfg.HostPort,
	}

	// Configure TLS if enabled
	if cfg.TLS.Enabled {
		tlsConfig, err := cfg.TLS.GetTLSConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to configure TLS: %w", err)
		}
		options.ConnectionOptions = client.ConnectionOptions{
			TLS: tlsConfig,
		}
	}

	return client.NewNamespaceClient(options)
}

// parseRetentionPeriod 解析保留期
func parseRetentionPeriod(retentionStr string) *durationpb.Duration {
	if retentionStr == "" {
		return durationpb.New(defaultRetentionPeriod)
	}

	if duration, err := time.ParseDuration(retentionStr); err == nil {
		return durationpb.New(duration)
	}

	return durationpb.New(defaultRetentionPeriod)
}

// getArchivalState 获取归档状态
func getArchivalState(enabled bool) enums.ArchivalState {
	if enabled {
		return enums.ARCHIVAL_STATE_ENABLED
	}
	return enums.ARCHIVAL_STATE_DISABLED
}

// isNamespaceNotFoundError 判断是否为命名空间不存在错误
func isNamespaceNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToLower(err.Error())
	return strings.Contains(errMsg, "namespace") &&
		(strings.Contains(errMsg, "not found") ||
			strings.Contains(errMsg, "does not exist") ||
			strings.Contains(errMsg, "unknown"))
}

// ValidateNamespace 验证命名空间是否存在
func ValidateNamespace(ctx context.Context, cfg *config.Config) error {
	nsClient, err := createNamespaceClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create namespace client: %w", err)
	}
	defer nsClient.Close()

	if _, err := nsClient.Describe(ctx, cfg.Namespace); err != nil {
		return fmt.Errorf("namespace '%s' validation failed: %w", cfg.Namespace, err)
	}

	return nil
}

// CreateMultipleNamespaces 批量创建命名空间
func CreateMultipleNamespaces(ctx context.Context, cfg *config.Config, namespaces []string) error {
	nsClient, err := createNamespaceClient(cfg)
	if err != nil {
		return fmt.Errorf("failed to create namespace client: %w", err)
	}
	defer nsClient.Close()

	for _, ns := range namespaces {
		nsCfg := *cfg
		nsCfg.Namespace = ns

		// 检查并创建
		if _, err = nsClient.Describe(ctx, ns); err != nil {
			if isNamespaceNotFoundError(err) {
				if err = createNamespace(ctx, nsClient, &nsCfg); err != nil {
					return fmt.Errorf("failed to create namespace '%s': %w", ns, err)
				}
			} else {
				return fmt.Errorf("failed to check namespace '%s': %w", ns, err)
			}
		}
	}

	return nil
}
