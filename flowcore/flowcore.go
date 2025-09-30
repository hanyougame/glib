package flowcore

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hanyougame/glib/flowcore/config"
	"github.com/hanyougame/glib/flowcore/core"
	"github.com/hanyougame/glib/flowcore/interceptor"
	"github.com/hanyougame/glib/flowcore/logger"
	flowcorepkg "github.com/hanyougame/glib/flowcore/pkg"
	enumspb "go.temporal.io/api/enums/v1"
	sdkclient "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// ========================= SDK 状态管理 =========================

var (
	sdk *SDK
	mu  sync.RWMutex
)

// SDK 包含所有 SDK 组件
type SDK struct {
	config      *config.Config
	initialized bool
}

// ========================= 初始化选项 =========================

// InitializeOptions SDK 初始化选项
type InitializeOptions struct {
	Config            *config.Config      // 基本配置
	MetricsConfig     *flowcorepkg.Config // 指标配置
	InterceptorConfig *interceptor.Config // 拦截器配置
}

// DefaultInitializeOptions 返回默认初始化选项
func DefaultInitializeOptions() *InitializeOptions {
	return &InitializeOptions{
		Config:            config.DefaultConfig(),
		MetricsConfig:     flowcorepkg.DefaultConfig(),
		InterceptorConfig: interceptor.DefaultConfig(),
	}
}

// ========================= SDK 初始化 =========================

// Initialize 初始化 Temporal SDK
func Initialize(ctx context.Context, cfg *config.Config) error {
	return InitializeWithOptions(ctx, &InitializeOptions{Config: cfg})
}

// InitializeWithOptions 使用选项初始化 Temporal SDK
func InitializeWithOptions(ctx context.Context, opts *InitializeOptions) error {
	mu.Lock()
	defer mu.Unlock()

	if sdk != nil && sdk.initialized {
		return fmt.Errorf("temporal SDK already initialized")
	}

	// 准备选项
	if opts == nil {
		opts = DefaultInitializeOptions()
	}
	if opts.Config == nil {
		opts.Config = config.DefaultConfig()
	}

	// 验证配置
	if err := validateConfig(opts.Config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// 初始化各个组件
	if err := initializeComponents(ctx, opts); err != nil {
		return err
	}

	// 创建 SDK 实例
	sdk = &SDK{
		config:      opts.Config,
		initialized: true,
	}

	fmt.Println("Temporal SDK initialized successfully")
	fmt.Printf("Workers configured: %d\n", len(opts.Config.Workers))

	return nil
}

// initializeComponents 初始化所有组件
func initializeComponents(ctx context.Context, opts *InitializeOptions) error {
	// 初始化重试策略
	flowcorepkg.InitializeDefaultPolicies()

	// 初始化指标管理器
	if err := initializeMetrics(opts); err != nil {
		return fmt.Errorf("failed to initialize metrics: %w", err)
	}

	// 创建默认客户端
	if _, err := core.GetOrCreateDefault(ctx, opts.Config); err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	return nil
}

// initializeMetrics 初始化指标系统
func initializeMetrics(opts *InitializeOptions) error {
	if opts.MetricsConfig == nil {
		opts.MetricsConfig = flowcorepkg.DefaultConfig()
	}

	// 设置日志器
	if opts.MetricsConfig.Logger == nil {
		opts.MetricsConfig.Logger = logger.NewDevelopmentLogger("info")
	}

	flowcorepkg.InitGlobalManager(opts.MetricsConfig)
	return nil
}

// MustInitialize 初始化 SDK，如果失败则 panic
func MustInitialize(ctx context.Context, cfg *config.Config) {
	if err := Initialize(ctx, cfg); err != nil {
		panic(fmt.Sprintf("Failed to initialize Temporal SDK: %v", err))
	}
}

// ========================= 客户端访问 =========================

// Client 返回默认的 Temporal 客户端实例
func Client(ctx context.Context) sdkclient.Client {
	ensureInitialized()

	c, err := core.GetOrCreateDefault(ctx, sdk.config)
	if err != nil {
		panic(fmt.Sprintf("failed to get client: %v", err))
	}
	return c
}

// ========================= 注册功能 =========================

// RegisterWorkflow 全局注册一个工作流
func RegisterWorkflow(name string, workflow interface{}) {
	core.RegisterWorkflow(name, workflow)
}

// RegisterActivity 全局注册一个活动
func RegisterActivity(name string, activity interface{}) {
	core.RegisterActivity(name, activity)
}

// RegisterModule 全局注册一个模块
func RegisterModule(module core.Module) {
	core.RegisterModule(module)
}

// ========================= Worker 管理 =========================

// WorkerOptions Worker 创建选项
type WorkerOptions struct {
	Config            *config.WorkerConfig // Worker 配置
	InterceptorConfig *interceptor.Config  // 拦截器配置
	OnError           func(error)          // 错误回调
}

// NewWorker 创建一个新的工作者
func NewWorker(ctx context.Context, cfg *config.WorkerConfig) *core.Worker {
	return NewWorkerWithOptions(ctx, &WorkerOptions{Config: cfg})
}

// NewWorkerWithOptions 使用选项创建一个新的工作者
func NewWorkerWithOptions(ctx context.Context, opts *WorkerOptions) *core.Worker {
	ensureInitialized()

	if opts == nil || opts.Config == nil {
		panic("worker config is required")
	}

	// 准备拦截器配置
	interceptorConfig := opts.InterceptorConfig
	if interceptorConfig == nil {
		interceptorConfig = interceptor.DefaultConfig()
	}

	// 创建 Worker 选项
	workerOpts := &core.WorkerOptions{
		EnableTrace:      interceptorConfig.EnableTrace,
		EnableMetrics:    interceptorConfig.EnableMetrics,
		Logger:           interceptorConfig.Logger,
		CustomAttributes: interceptorConfig.CustomAttributes,
		OnError:          opts.OnError,
	}

	// 创建工作者
	c := Client(ctx)
	w := core.NewWorkerWithOptions(c, opts.Config, workerOpts)

	// 应用全局注册表中的工作流和活动
	core.GetGlobal().ApplyTo(w)

	// 添加到管理器
	if err := core.AddWorker(w); err != nil {
		// 如果添加失败（重复的 TaskQueue），只记录警告
		fmt.Printf("Warning: %v\n", err)
	}

	return w
}

// StartWorkers 启动所有工作者
func StartWorkers(ctx context.Context) error {
	ensureInitialized()
	return core.StartAll(ctx)
}

// StopWorkers 停止所有工作者
func StopWorkers() {
	core.StopAll()
}

// ========================= 工作流和活动选项 =========================

// WorkflowOptions 创建标准的工作流选项
func WorkflowOptions(workflowID, taskQueue string) sdkclient.StartWorkflowOptions {
	ensureInitialized()

	return sdkclient.StartWorkflowOptions{
		ID:                       workflowID,
		TaskQueue:                taskQueue,
		WorkflowExecutionTimeout: sdk.config.Timeout.WorkflowExecution,
		WorkflowTaskTimeout:      sdk.config.Timeout.WorkflowTask,
		WorkflowIDReusePolicy:    enumspb.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE,
	}
}

// WorkflowOptionsWithRetry 创建带有自定义重试策略的工作流选项
func WorkflowOptionsWithRetry(workflowID, taskQueue string, retryPolicy *temporal.RetryPolicy) sdkclient.StartWorkflowOptions {
	options := WorkflowOptions(workflowID, taskQueue)
	options.RetryPolicy = retryPolicy
	return options
}

// ActivityOptions 创建标准的活动选项
func ActivityOptions(options ...core.ActivityOption) workflow.ActivityOptions {
	ensureInitialized()

	opts := workflow.ActivityOptions{
		ScheduleToStartTimeout: sdk.config.Timeout.ActivityStart,
		StartToCloseTimeout:    sdk.config.Timeout.ActivityExecution,
		HeartbeatTimeout:       sdk.config.Timeout.ActivityHeartbeat,
		RetryPolicy:            flowcorepkg.Get(flowcorepkg.Standard),
	}

	// 应用自定义选项
	for _, option := range options {
		option(&opts)
	}

	return opts
}

// ActivityOptionsWithRetry 创建带有自定义重试策略的活动选项
func ActivityOptionsWithRetry(retryPolicyName string, options ...core.ActivityOption) workflow.ActivityOptions {
	opts := ActivityOptions(options...)
	opts.RetryPolicy = flowcorepkg.Get(retryPolicyName)
	return opts
}

// CriticalActivityOptions 创建关键活动的活动选项
func CriticalActivityOptions(options ...core.ActivityOption) workflow.ActivityOptions {
	return ActivityOptionsWithRetry(flowcorepkg.Critical, options...)
}

// LocalActivityOptions 创建标准的本地活动选项
func LocalActivityOptions() workflow.LocalActivityOptions {
	ensureInitialized()

	return workflow.LocalActivityOptions{
		ScheduleToCloseTimeout: sdk.config.Timeout.ActivityExecution,
		RetryPolicy:            flowcorepkg.Get(flowcorepkg.Standard),
	}
}

// ========================= 预设选项 =========================

// HighFrequencyFastOptions 高频快速处理选项
func HighFrequencyFastOptions() workflow.ActivityOptions {
	return ActivityOptions(
		core.WithScheduleToStartTimeout(5*time.Second),
		core.WithStartToCloseTimeout(30*time.Second),
		core.WithHeartbeatTimeout(10*time.Second),
		core.WithRetryPolicy(flowcorepkg.Critical),
	)
}

// LowFrequencySlowOptions 低频慢速处理选项
func LowFrequencySlowOptions() workflow.ActivityOptions {
	return ActivityOptions(
		core.WithScheduleToStartTimeout(2*time.Minute),
		core.WithStartToCloseTimeout(60*time.Minute),
		core.WithHeartbeatTimeout(2*time.Minute),
		core.WithRetryPolicy(flowcorepkg.Critical),
	)
}

// NormalProcessingOptions 正常处理选项
func NormalProcessingOptions() workflow.ActivityOptions {
	return ActivityOptions(
		core.WithScheduleToStartTimeout(30*time.Second),
		core.WithStartToCloseTimeout(5*time.Minute),
		core.WithHeartbeatTimeout(30*time.Second),
		core.WithRetryPolicy(flowcorepkg.Critical),
	)
}

// ========================= SDK 生命周期 =========================

// Shutdown 关闭 SDK
func Shutdown() {
	mu.Lock()
	defer mu.Unlock()

	if sdk == nil || !sdk.initialized {
		return
	}

	// 停止所有 Workers
	core.StopAll()

	// 停止指标管理器
	flowcorepkg.StopGlobalManager()

	// 关闭所有客户端
	core.CloseAll()

	sdk.initialized = false
}

// IsInitialized 返回 SDK 是否已初始化
func IsInitialized() bool {
	mu.RLock()
	defer mu.RUnlock()
	return sdk != nil && sdk.initialized
}

// ========================= 信息和指标 =========================

// GetConfig 返回当前配置
func GetConfig() *config.Config {
	mu.RLock()
	defer mu.RUnlock()

	if sdk == nil {
		return nil
	}
	return sdk.config
}

// GetStats 返回 SDK 统计信息
func GetStats() map[string]interface{} {
	stats := map[string]interface{}{
		"initialized": IsInitialized(),
	}

	if !IsInitialized() {
		return stats
	}

	// 获取注册信息
	workflowCount, activityCount := core.GetGlobal().Count()
	stats["workflow_count"] = workflowCount
	stats["activity_count"] = activityCount
	stats["worker_count"] = core.GetWorkerCount()
	stats["workers_running"] = core.IsRunning()

	// 添加指标统计
	for k, v := range flowcorepkg.GetGlobalMetrics() {
		stats[fmt.Sprintf("metrics_%s", k)] = v
	}

	return stats
}

// ========================= 辅助函数 =========================

// ensureInitialized 确保 SDK 已初始化
func ensureInitialized() {
	mu.RLock()
	defer mu.RUnlock()

	if sdk == nil || !sdk.initialized {
		panic("temporal SDK not initialized - call Initialize() first")
	}
}

// validateConfig 验证配置的有效性
func validateConfig(cfg *config.Config) error {
	if cfg.HostPort == "" {
		return fmt.Errorf("host_port is required")
	}
	if cfg.Namespace == "" {
		return fmt.Errorf("namespace is required")
	}
	if len(cfg.Workers) == 0 {
		return fmt.Errorf("at least one worker configuration is required")
	}
	for i, workerCfg := range cfg.Workers {
		if workerCfg.TaskQueue == "" {
			return fmt.Errorf("worker[%d] task_queue is required", i)
		}
	}
	return nil
}
