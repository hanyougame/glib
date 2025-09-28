package config

import (
	"crypto/tls"
	"time"
)

// Config Temporal 配置
type Config struct {
	// 连接设置
	HostPort  string `json:"host_port,optional" yaml:"host_port"`
	Namespace string `json:"namespace,default=default_workflow" yaml:"namespace"`
	Identity  string `json:"identity,optional" yaml:"identity"`

	// TLS 配置
	TLS TLSConfig `json:"tls,optional" yaml:"tls"`

	// Worker 设置
	Workers []WorkerConfig `json:"workers,optional" yaml:"workers"`

	// 重试和超时设置
	Retry   RetryConfig   `json:"retry,optional" yaml:"retry"`
	Timeout TimeoutConfig `json:"timeout,optional" yaml:"timeout"`

	// 高级设置
	DataConverter string `json:"data_converter,default=json" yaml:"data_converter"`
	Enabled       bool   `json:"enabled,default=true" yaml:"enabled"`
	// 命名空间配置
	NamespaceConf NamespaceConfig `json:"namespace_conf,optional" yaml:"namespace_conf"`
}

type TLSConfig struct {
	Enabled            bool   `json:"enabled,optional" yaml:"enabled"`
	CertFile           string `json:"cert_file,optional" yaml:"cert_file"`
	KeyFile            string `json:"key_file,optional" yaml:"key_file"`
	CAFile             string `json:"ca_file,optional" yaml:"ca_file"`
	ServerName         string `json:"server_name,optional" yaml:"server_name"`
	InsecureSkipVerify bool   `json:"insecure_skip_verify,optional" yaml:"insecure_skip_verify"`
}

type WorkerConfig struct {
	TaskQueue               string        `json:"task_queue,optional" yaml:"task_queue"`
	MaxConcurrentActivities int           `json:"max_concurrent_activities,optional" yaml:"max_concurrent_activities"`
	MaxConcurrentWorkflows  int           `json:"max_concurrent_workflows,optional" yaml:"max_concurrent_workflows"`
	EnableSessionWorker     bool          `json:"enable_session_worker,default=false" yaml:"enable_session_worker"`
	StickyScheduleToStart   time.Duration `json:"sticky_schedule_to_start,optional" yaml:"sticky_schedule_to_start"`
}

type RetryConfig struct {
	InitialInterval    time.Duration `json:"initial_interval,optional" yaml:"initial_interval"`
	BackoffCoefficient float64       `json:"backoff_coefficient,optional" yaml:"backoff_coefficient"`
	MaximumInterval    time.Duration `json:"maximum_interval,optional" yaml:"maximum_interval"`
	MaximumAttempts    int32         `json:"maximum_attempts,optional" yaml:"maximum_attempts"`
}

type TimeoutConfig struct {
	WorkflowExecution time.Duration `json:"workflow_execution,optional" yaml:"workflow_execution"`
	WorkflowTask      time.Duration `json:"workflow_task,optional" yaml:"workflow_task"`
	ActivityExecution time.Duration `json:"activity_execution,optional" yaml:"activity_execution"`
	ActivityStart     time.Duration `json:"activity_start,optional" yaml:"activity_start"`
	ActivityHeartbeat time.Duration `json:"activity_heartbeat,optional" yaml:"activity_heartbeat"`
}

type NamespaceConfig struct {
	AutoCreate                bool   `json:"auto_create,default=true" yaml:"auto_create"`                                 // 是否自动创建命名空间
	Description               string `json:"description,optional" yaml:"description"`                                     // 命名空间描述
	RetentionPeriod           string `json:"retention_period,default=72h" yaml:"retention_period"`                        // 保留期 (如 "72h", "30d")
	HistoryArchivalEnabled    bool   `json:"history_archival_enabled,default=true" yaml:"history_archival_enabled"`       // 是否启用历史归档
	VisibilityArchivalEnabled bool   `json:"visibility_archival_enabled,default=true" yaml:"visibility_archival_enabled"` // 是否启用可见性归档
}

// GetTLSConfig 返回配置的 TLS 配置
func (c *TLSConfig) GetTLSConfig() (*tls.Config, error) {
	if !c.Enabled {
		return nil, nil
	}

	tlsConfig := &tls.Config{
		ServerName:         c.ServerName,
		InsecureSkipVerify: c.InsecureSkipVerify,
	}

	if c.CertFile != "" && c.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	return tlsConfig, nil
}

// DefaultConfig 返回一个默认配置
func DefaultConfig() *Config {
	return &Config{
		HostPort:  "localhost:7233",
		Namespace: "default",
		Enabled:   true,
		Workers: []WorkerConfig{
			{
				TaskQueue:               "default",
				MaxConcurrentActivities: 100,
				MaxConcurrentWorkflows:  100,
				StickyScheduleToStart:   10 * time.Second,
			},
		},
		Retry: RetryConfig{
			InitialInterval:    time.Second,
			BackoffCoefficient: 2.0,
			MaximumInterval:    time.Minute,
			MaximumAttempts:    5,
		},
		Timeout: TimeoutConfig{
			WorkflowExecution: 24 * time.Hour,
			WorkflowTask:      10 * time.Second,
			ActivityExecution: 30 * time.Minute,
			ActivityStart:     5 * time.Second,
			ActivityHeartbeat: 10 * time.Second,
		},
		DataConverter: "json",
	}
}
