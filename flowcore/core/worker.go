package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hanyougame/glib/flowcore/config"
	flowcoreinterceptor "github.com/hanyougame/glib/flowcore/interceptor"
	flowcorepkg "github.com/hanyougame/glib/flowcore/pkg"
	sdkinterceptor "go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/log"

	sdkactivity "go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	sdkworkflow "go.temporal.io/sdk/workflow"
)

// ========================= Worker 状态定义 =========================

// WorkerStatus 工作器状态
type WorkerStatus int

const (
	WorkerStatusIdle    WorkerStatus = iota // 空闲
	WorkerStatusRunning                     // 运行中
	WorkerStatusStopped                     // 已停止
	WorkerStatusError                       // 错误
)

// String 返回状态的字符串表示
func (s WorkerStatus) String() string {
	switch s {
	case WorkerStatusIdle:
		return "idle"
	case WorkerStatusRunning:
		return "running"
	case WorkerStatusStopped:
		return "stopped"
	case WorkerStatusError:
		return "error"
	default:
		return "unknown"
	}
}

// ========================= Worker 结构 =========================

// Worker 包含了一个 Temporal worker
type Worker struct {
	worker      worker.Worker
	config      *config.WorkerConfig
	client      client.Client
	interceptor *flowcoreinterceptor.UnifiedInterceptor
	startTime   time.Time
	status      WorkerStatus
	mu          sync.RWMutex
}

// WorkerOptions worker 配置选项
type WorkerOptions struct {
	EnableTrace      bool              // 启用 trace
	EnableMetrics    bool              // 启用指标收集
	Logger           log.Logger        // 自定义日志器
	CustomAttributes map[string]string // 自定义属性
	OnError          func(error)       // 错误回调
}

// DefaultWorkerOptions 默认配置
func DefaultWorkerOptions() *WorkerOptions {
	return &WorkerOptions{
		EnableTrace:      true,
		EnableMetrics:    true,
		CustomAttributes: make(map[string]string),
	}
}

// ========================= Worker 创建 =========================

// NewWorker 创建一个新的 worker（使用默认选项）
func NewWorker(c client.Client, cfg *config.WorkerConfig) *Worker {
	return NewWorkerWithOptions(c, cfg, DefaultWorkerOptions())
}

// NewWorkerWithOptions 创建带自定义选项的 worker
func NewWorkerWithOptions(c client.Client, cfg *config.WorkerConfig, opts *WorkerOptions) *Worker {
	if opts == nil {
		opts = DefaultWorkerOptions()
	}

	// 创建统一拦截器配置
	interceptorConfig := &flowcoreinterceptor.Config{
		EnableTrace:      opts.EnableTrace,
		EnableMetrics:    opts.EnableMetrics,
		Logger:           opts.Logger,
		CustomAttributes: opts.CustomAttributes,
	}

	// 创建 Worker
	return createWorker(c, cfg, interceptorConfig, opts)
}

// createWorker 内部创建 worker 的方法
func createWorker(c client.Client, cfg *config.WorkerConfig, interceptorConfig *flowcoreinterceptor.Config, opts *WorkerOptions) *Worker {

	// 创建统一拦截器
	unifiedInterceptor := flowcoreinterceptor.NewUnifiedInterceptor(interceptorConfig).(*flowcoreinterceptor.UnifiedInterceptor)

	// 构建 Worker 选项
	workerOptions := worker.Options{
		MaxConcurrentActivityExecutionSize:     cfg.MaxConcurrentActivities,
		MaxConcurrentWorkflowTaskExecutionSize: cfg.MaxConcurrentWorkflows,
		EnableSessionWorker:                    cfg.EnableSessionWorker,
		StickyScheduleToStartTimeout:           cfg.StickyScheduleToStart,
		Interceptors: []sdkinterceptor.WorkerInterceptor{
			unifiedInterceptor,
		},
	}

	// 设置错误回调
	if opts != nil && opts.OnError != nil {
		workerOptions.OnFatalError = opts.OnError
	}

	// 创建 SDK Worker
	w := worker.New(c, cfg.TaskQueue, workerOptions)

	return &Worker{
		worker:      w,
		config:      cfg,
		client:      c,
		interceptor: unifiedInterceptor,
		status:      WorkerStatusIdle,
	}
}

// ========================= Worker 注册方法 =========================

// RegisterWorkflow 注册一个工作流
func (w *Worker) RegisterWorkflow(workflow interface{}) {
	w.worker.RegisterWorkflow(workflow)
}

// RegisterWorkflowWithOptions 注册带选项的工作流
func (w *Worker) RegisterWorkflowWithOptions(workflow interface{}, options sdkworkflow.RegisterOptions) {
	w.worker.RegisterWorkflowWithOptions(workflow, options)
}

// RegisterActivity 注册一个活动
func (w *Worker) RegisterActivity(activity interface{}) {
	w.worker.RegisterActivity(activity)
}

// RegisterActivityWithOptions 注册带选项的活动
func (w *Worker) RegisterActivityWithOptions(activity interface{}, options sdkactivity.RegisterOptions) {
	w.worker.RegisterActivityWithOptions(activity, options)
}

// ========================= Worker 生命周期管理 =========================

// Start 启动 worker
func (w *Worker) Start(ctx context.Context) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.status == WorkerStatusRunning {
		return fmt.Errorf("worker already running for task queue: %s", w.config.TaskQueue)
	}

	// 在 goroutine 中运行 worker
	go w.run()

	return nil
}

// run 内部运行方法
func (w *Worker) run() {
	w.setStatus(WorkerStatusRunning)
	w.startTime = time.Now()

	fmt.Printf("Starting worker for task queue: %s\n", w.config.TaskQueue)
	err := w.worker.Run(worker.InterruptCh())

	if err != nil {
		w.setStatus(WorkerStatusError)
		fmt.Printf("Worker error for task queue %s: %v\n", w.config.TaskQueue, err)
	} else {
		w.setStatus(WorkerStatusStopped)
	}
}

// Stop 停止 worker
func (w *Worker) Stop() {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.status != WorkerStatusRunning {
		return
	}

	w.worker.Stop()
	w.status = WorkerStatusStopped

	fmt.Printf("Worker stopped for task queue: %s\n", w.config.TaskQueue)
}

// setStatus 线程安全地设置状态
func (w *Worker) setStatus(status WorkerStatus) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.status = status
}

// ========================= Worker 信息获取 =========================

// GetTaskQueue 返回任务队列名称
func (w *Worker) GetTaskQueue() string {
	return w.config.TaskQueue
}

// GetStatus 获取工作器状态
func (w *Worker) GetStatus() WorkerStatus {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.status
}

// GetUptime 获取运行时间
func (w *Worker) GetUptime() time.Duration {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if w.status != WorkerStatusRunning {
		return 0
	}
	return time.Since(w.startTime)
}

// GetMetrics 获取当前 worker 的指标
func (w *Worker) GetMetrics() map[string]int64 {
	if w.interceptor != nil {
		return w.interceptor.GetMetrics()
	}
	return nil
}

// GetDetailedMetrics 获取详细指标
func (w *Worker) GetDetailedMetrics() (map[string]*flowcorepkg.ExecutionMetric, map[string]*flowcorepkg.ExecutionMetric) {
	if w.interceptor != nil {
		return w.interceptor.GetDetailedMetrics()
	}
	return nil, nil

}

// GetInfo 获取 Worker 信息
func (w *Worker) GetInfo() WorkerInfo {
	w.mu.RLock()
	defer w.mu.RUnlock()

	info := WorkerInfo{
		TaskQueue:               w.config.TaskQueue,
		Status:                  w.status,
		StartTime:               w.startTime,
		Uptime:                  w.GetUptime(),
		MaxConcurrentActivities: w.config.MaxConcurrentActivities,
		MaxConcurrentWorkflows:  w.config.MaxConcurrentWorkflows,
	}

	// 获取指标信息
	if metrics := w.GetMetrics(); metrics != nil {
		info.ActivityCount = metrics["activities"]
		info.WorkflowCount = metrics["workflows"]
		info.ErrorCount = metrics["errors"]
	}

	return info
}

// WorkerInfo Worker 信息
type WorkerInfo struct {
	TaskQueue               string
	Status                  WorkerStatus
	StartTime               time.Time
	Uptime                  time.Duration
	MaxConcurrentActivities int
	MaxConcurrentWorkflows  int
	ActivityCount           int64
	WorkflowCount           int64
	ErrorCount              int64
}

// ========================= Worker 管理器 =========================

// Manager 管理多个 workers
type Manager struct {
	mu             sync.RWMutex
	workers        []*Worker
	workersByQueue map[string]*Worker // 按任务队列索引
	running        bool
}

// 全局管理器实例
var globalManager = &Manager{
	workers:        make([]*Worker, 0),
	workersByQueue: make(map[string]*Worker),
}

// AddWorker 将 worker 添加到全局管理器
func AddWorker(w *Worker) error {
	globalManager.mu.Lock()
	defer globalManager.mu.Unlock()

	// 检查是否已存在相同任务队列的 worker
	if _, exists := globalManager.workersByQueue[w.config.TaskQueue]; exists {
		return fmt.Errorf("worker already exists for task queue: %s", w.config.TaskQueue)
	}

	globalManager.workers = append(globalManager.workers, w)
	globalManager.workersByQueue[w.config.TaskQueue] = w

	fmt.Printf("Added worker for task queue: %s\n", w.config.TaskQueue)
	return nil
}

// RemoveWorker 从管理器中移除 worker
func RemoveWorker(taskQueue string) error {
	globalManager.mu.Lock()
	defer globalManager.mu.Unlock()

	_worker, exists := globalManager.workersByQueue[taskQueue]
	if !exists {
		return fmt.Errorf("worker not found for task queue: %s", taskQueue)
	}

	// 停止 worker
	_worker.Stop()

	// 从 workers 列表中移除
	for i, w := range globalManager.workers {
		if w.config.TaskQueue == taskQueue {
			globalManager.workers = append(globalManager.workers[:i], globalManager.workers[i+1:]...)
			break
		}
	}

	// 从索引中移除
	delete(globalManager.workersByQueue, taskQueue)

	fmt.Printf("Removed worker for task queue: %s\n", taskQueue)
	return nil
}

// StartAll 启动所有 workers
func StartAll(ctx context.Context) error {
	globalManager.mu.Lock()
	defer globalManager.mu.Unlock()

	if globalManager.running {
		return fmt.Errorf("workers already running")
	}

	return startWorkers(ctx)
}

// startWorkers 内部启动 workers 的方法
func startWorkers(ctx context.Context) error {
	fmt.Printf("Starting %d workers\n", len(globalManager.workers))

	var firstErr error
	successCount := 0

	for _, w := range globalManager.workers {
		if err := w.Start(ctx); err != nil {
			if firstErr == nil {
				firstErr = err
			}
			fmt.Printf("Failed to start worker for task queue %s: %v\n", w.config.TaskQueue, err)
		} else {
			successCount++
		}
	}

	if successCount > 0 {
		globalManager.running = true
	}

	if firstErr != nil {
		return fmt.Errorf("started %d/%d workers, first error: %w",
			successCount, len(globalManager.workers), firstErr)
	}

	return nil
}

// StopAll 停止所有 workers
func StopAll() {
	globalManager.mu.Lock()
	defer globalManager.mu.Unlock()

	if !globalManager.running {
		return
	}

	stopWorkers()
}

// stopWorkers 内部停止 workers 的方法
func stopWorkers() {
	fmt.Printf("Stopping %d workers\n", len(globalManager.workers))

	for _, w := range globalManager.workers {
		w.Stop()
	}

	globalManager.running = false
	fmt.Println("All workers stopped")
}

// ========================= Worker 查询方法 =========================

// GetWorkerCount 返回已注册的 worker 数量
func GetWorkerCount() int {
	globalManager.mu.RLock()
	defer globalManager.mu.RUnlock()
	return len(globalManager.workers)
}

// GetRunningWorkerCount 返回运行中的 worker 数量
func GetRunningWorkerCount() int {
	globalManager.mu.RLock()
	defer globalManager.mu.RUnlock()

	count := 0
	for _, w := range globalManager.workers {
		if w.GetStatus() == WorkerStatusRunning {
			count++
		}
	}
	return count
}

// IsRunning 返回 workers 是否正在运行
func IsRunning() bool {
	globalManager.mu.RLock()
	defer globalManager.mu.RUnlock()
	return globalManager.running
}

// GetWorkerByTaskQueue 根据任务队列名称获取 worker
func GetWorkerByTaskQueue(taskQueue string) *Worker {
	globalManager.mu.RLock()
	defer globalManager.mu.RUnlock()
	return globalManager.workersByQueue[taskQueue]
}

// GetWorkersInfo 获取所有 workers 的信息
func GetWorkersInfo() []WorkerInfo {
	globalManager.mu.RLock()
	defer globalManager.mu.RUnlock()

	infos := make([]WorkerInfo, 0, len(globalManager.workers))
	for _, w := range globalManager.workers {
		infos = append(infos, w.GetInfo())
	}
	return infos
}

// ========================= 指标相关方法 =========================

// GetAllMetrics 获取所有 workers 的合并指标
func GetAllMetrics() map[string]int64 {
	globalManager.mu.RLock()
	defer globalManager.mu.RUnlock()

	totalMetrics := map[string]int64{
		"total_workers":   int64(len(globalManager.workers)),
		"running_workers": 0,
	}

	// 计算运行中的 workers
	for _, w := range globalManager.workers {
		if w.GetStatus() == WorkerStatusRunning {
			totalMetrics["running_workers"]++
		}
	}

	// 获取并合并指标（使用全局收集器，所以直接获取一次即可）
	if len(globalManager.workers) > 0 {
		if metrics := globalManager.workers[0].GetMetrics(); metrics != nil {
			for key, value := range metrics {
				totalMetrics[key] = value
			}
		}
	}

	return totalMetrics
}
