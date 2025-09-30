package pkg

import (
	"sync"
	"sync/atomic"
	"time"

	"go.temporal.io/sdk/log"
)

// ========================= 配置 =========================

type Config struct {
	// 上报配置
	ReportInterval   time.Duration
	EnableAutoReport bool
	CustomReporter   func(metrics map[string]interface{})

	// 滑动窗口配置
	MetricsWindowSize     int           // 滑动窗口大小（保留多少个时间段的数据）
	MetricsWindowDuration time.Duration // 每个窗口的时长

	// 重置配置
	EnableAutoReset bool          // 是否自动重置累积指标
	ResetInterval   time.Duration // 重置间隔

	// 详细指标
	EnableDetailedMetrics bool

	// 日志器
	Logger log.Logger
}

func DefaultConfig() *Config {
	return &Config{
		ReportInterval:        5 * time.Minute,
		MetricsWindowSize:     24,               // 保留24个窗口
		MetricsWindowDuration: 10 * time.Minute, // 每个窗口10分钟
		EnableAutoReport:      true,
		EnableAutoReset:       true,
		ResetInterval:         24 * time.Hour, // 每24小时重置一次累积值
		EnableDetailedMetrics: false,
	}
}

// ========================= 执行指标 =========================

// ExecutionMetric 执行指标（用于详细指标）
type ExecutionMetric struct {
	Name          string
	Count         int64
	SuccessCount  int64
	FailureCount  int64
	TotalDuration int64
	MaxDuration   int64
	MinDuration   int64
	RetryCount    int64
	LastError     string
	LastExecution time.Time
	mu            sync.RWMutex // 保护 LastError 和 LastExecution
}

// ========================= 滑动窗口指标 =========================

// MetricsWindow 单个时间窗口的指标
type MetricsWindow struct {
	StartTime            time.Time
	EndTime              time.Time
	ActivityCount        int64
	WorkflowCount        int64
	ActivitySuccessCount int64
	ActivityFailureCount int64
	WorkflowSuccessCount int64
	WorkflowFailureCount int64
	ActivityDuration     int64
	WorkflowDuration     int64
	MaxActivityDuration  int64
	MaxWorkflowDuration  int64
	ErrorCount           int64
	ActivityRetryCount   int64
	WorkflowRetryCount   int64
}

// SlidingWindowMetrics 滑动窗口指标管理
type SlidingWindowMetrics struct {
	windows      []*MetricsWindow
	currentIndex int
	windowSize   int
	duration     time.Duration
	mu           sync.RWMutex
	lastRotate   time.Time
}

func NewSlidingWindowMetrics(size int, duration time.Duration) *SlidingWindowMetrics {
	windows := make([]*MetricsWindow, size)
	now := time.Now()
	for i := 0; i < size; i++ {
		windows[i] = &MetricsWindow{
			StartTime: now,
			EndTime:   now.Add(duration),
		}
	}

	return &SlidingWindowMetrics{
		windows:      windows,
		windowSize:   size,
		duration:     duration,
		lastRotate:   now,
		currentIndex: 0,
	}
}

// rotateIfNeeded 检查并旋转窗口
func (s *SlidingWindowMetrics) rotateIfNeeded() {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(s.lastRotate)

	// 计算需要旋转多少个窗口
	rotations := int(elapsed / s.duration)
	if rotations <= 0 {
		return
	}

	// 限制旋转数量，避免一次旋转太多
	if rotations > s.windowSize {
		rotations = s.windowSize
	}

	for i := 0; i < rotations; i++ {
		// 移动到下一个窗口
		s.currentIndex = (s.currentIndex + 1) % s.windowSize
		// 重置新窗口
		s.windows[s.currentIndex] = &MetricsWindow{
			StartTime: now.Add(time.Duration(i-rotations+1) * s.duration),
			EndTime:   now.Add(time.Duration(i-rotations+2) * s.duration),
		}
	}

	s.lastRotate = now
}

// getCurrentWindow 获取当前窗口
func (s *SlidingWindowMetrics) getCurrentWindow() *MetricsWindow {
	s.rotateIfNeeded()
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.windows[s.currentIndex]
}

// getRecentMetrics 获取最近N个窗口的汇总指标
func (s *SlidingWindowMetrics) getRecentMetrics(count int) map[string]int64 {
	s.rotateIfNeeded()
	s.mu.RLock()
	defer s.mu.RUnlock()

	if count > s.windowSize {
		count = s.windowSize
	}

	result := make(map[string]int64)

	for i := 0; i < count; i++ {
		idx := (s.currentIndex - i + s.windowSize) % s.windowSize
		w := s.windows[idx]

		// 跳过未来的窗口
		if w.StartTime.After(time.Now()) {
			continue
		}

		result["activities"] += atomic.LoadInt64(&w.ActivityCount)
		result["workflows"] += atomic.LoadInt64(&w.WorkflowCount)
		result["errors"] += atomic.LoadInt64(&w.ErrorCount)
		result["activity_success"] += atomic.LoadInt64(&w.ActivitySuccessCount)
		result["activity_failure"] += atomic.LoadInt64(&w.ActivityFailureCount)
		result["workflow_success"] += atomic.LoadInt64(&w.WorkflowSuccessCount)
		result["workflow_failure"] += atomic.LoadInt64(&w.WorkflowFailureCount)
		result["activity_duration_ms"] += atomic.LoadInt64(&w.ActivityDuration)
		result["workflow_duration_ms"] += atomic.LoadInt64(&w.WorkflowDuration)
		result["activity_retries"] += atomic.LoadInt64(&w.ActivityRetryCount)
		result["workflow_retries"] += atomic.LoadInt64(&w.WorkflowRetryCount)

		// 最大值取所有窗口中的最大
		if max := atomic.LoadInt64(&w.MaxActivityDuration); max > result["max_activity_duration"] {
			result["max_activity_duration"] = max
		}
		if max := atomic.LoadInt64(&w.MaxWorkflowDuration); max > result["max_workflow_duration"] {
			result["max_workflow_duration"] = max
		}
	}

	return result
}

// ========================= 指标定义 =========================

// Metrics 统一的指标结构
type Metrics struct {
	// 累积指标（会定期重置）
	TotalActivityCount int64
	TotalWorkflowCount int64
	TotalErrorCount    int64
	ActivityRetryCount int64
	WorkflowRetryCount int64
	LastResetTime      time.Time

	// 当前并发（实时值，不累积）
	ConcurrentActivities int64
	ConcurrentWorkflows  int64

	// 滑动窗口指标
	slidingWindow *SlidingWindowMetrics

	// 详细指标 (可选)
	DetailedActivity sync.Map // map[string]*ExecutionMetric
	DetailedWorkflow sync.Map // map[string]*ExecutionMetric
}

// ========================= 指标收集器接口 =========================

type Collector interface {
	// Activity 相关
	RecordActivityStart(name string, attempt int32)
	RecordActivityEnd(name string, duration int64, err error)

	// Workflow 相关
	RecordWorkflowStart(name string)
	RecordWorkflowEnd(name string, duration int64, err error)

	// 并发控制
	IncrementConcurrentActivities()
	DecrementConcurrentActivities()
	IncrementConcurrentWorkflows()
	DecrementConcurrentWorkflows()

	// 数据获取
	GetMetrics() map[string]int64
	GetDetailedMetrics() (activities, workflows map[string]*ExecutionMetric)

	// 重置
	Reset()
}

// ========================= 标准收集器实现 =========================

type StandardCollector struct {
	metrics     *Metrics
	config      *Config
	mu          sync.RWMutex
	resetTicker *time.Ticker
	stopChan    chan struct{}
}

func NewStandardCollector(config *Config) *StandardCollector {
	if config == nil {
		config = DefaultConfig()
	}

	c := &StandardCollector{
		metrics: &Metrics{
			LastResetTime: time.Now(),
			slidingWindow: NewSlidingWindowMetrics(
				config.MetricsWindowSize,
				config.MetricsWindowDuration,
			),
		},
		config:   config,
		stopChan: make(chan struct{}),
	}

	// 启动自动重置
	if config.EnableAutoReset && config.ResetInterval > 0 {
		c.startAutoReset()
	}

	return c
}

// startAutoReset 启动自动重置
func (c *StandardCollector) startAutoReset() {
	c.resetTicker = time.NewTicker(c.config.ResetInterval)

	go func() {
		for {
			select {
			case <-c.resetTicker.C:
				c.resetAccumulatedMetrics()
			case <-c.stopChan:
				return
			}
		}
	}()
}

// resetAccumulatedMetrics 重置累积指标
func (c *StandardCollector) resetAccumulatedMetrics() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 重置累积计数
	atomic.StoreInt64(&c.metrics.TotalActivityCount, 0)
	atomic.StoreInt64(&c.metrics.TotalWorkflowCount, 0)
	atomic.StoreInt64(&c.metrics.TotalErrorCount, 0)
	atomic.StoreInt64(&c.metrics.ActivityRetryCount, 0)
	atomic.StoreInt64(&c.metrics.WorkflowRetryCount, 0)
	c.metrics.LastResetTime = time.Now()

	// 清空详细指标
	if c.config.EnableDetailedMetrics {
		c.metrics.DetailedActivity = sync.Map{}
		c.metrics.DetailedWorkflow = sync.Map{}
	}

	if c.config.Logger != nil {
		c.config.Logger.Info("Metrics reset completed", "reset_time", c.metrics.LastResetTime)
	}
}

// RecordActivityStart 记录活动开始
func (c *StandardCollector) RecordActivityStart(name string, attempt int32) {
	// 累积总数
	atomic.AddInt64(&c.metrics.TotalActivityCount, 1)

	// 当前窗口计数
	window := c.metrics.slidingWindow.getCurrentWindow()
	atomic.AddInt64(&window.ActivityCount, 1)

	if attempt > 1 {
		atomic.AddInt64(&c.metrics.ActivityRetryCount, 1)
		atomic.AddInt64(&window.ActivityRetryCount, 1)
	}

	// 详细指标
	if c.config.EnableDetailedMetrics {
		c.updateActivityMetricStart(name, attempt)
	}
}

// RecordActivityEnd 记录活动结束
func (c *StandardCollector) RecordActivityEnd(name string, duration int64, err error) {
	window := c.metrics.slidingWindow.getCurrentWindow()

	// 更新窗口内的持续时间
	atomic.AddInt64(&window.ActivityDuration, duration)

	// 更新最大执行时间
	c.updateMaxDuration(&window.MaxActivityDuration, duration)

	// 更新成功/失败计数
	if err != nil {
		atomic.AddInt64(&c.metrics.TotalErrorCount, 1)
		atomic.AddInt64(&window.ErrorCount, 1)
		atomic.AddInt64(&window.ActivityFailureCount, 1)
	} else {
		atomic.AddInt64(&window.ActivitySuccessCount, 1)
	}

	// 详细指标
	if c.config.EnableDetailedMetrics {
		c.updateActivityMetricEnd(name, duration, err)
	}
}

// RecordWorkflowStart 记录工作流开始
func (c *StandardCollector) RecordWorkflowStart(name string) {
	atomic.AddInt64(&c.metrics.TotalWorkflowCount, 1)

	window := c.metrics.slidingWindow.getCurrentWindow()
	atomic.AddInt64(&window.WorkflowCount, 1)

	// 详细指标
	if c.config.EnableDetailedMetrics {
		c.updateWorkflowMetricStart(name)
	}
}

// RecordWorkflowEnd 记录工作流结束
func (c *StandardCollector) RecordWorkflowEnd(name string, duration int64, err error) {
	window := c.metrics.slidingWindow.getCurrentWindow()

	atomic.AddInt64(&window.WorkflowDuration, duration)

	// 更新最大执行时间
	c.updateMaxDuration(&window.MaxWorkflowDuration, duration)

	if err != nil {
		atomic.AddInt64(&c.metrics.TotalErrorCount, 1)
		atomic.AddInt64(&window.ErrorCount, 1)
		atomic.AddInt64(&window.WorkflowFailureCount, 1)
	} else {
		atomic.AddInt64(&window.WorkflowSuccessCount, 1)
	}

	// 详细指标
	if c.config.EnableDetailedMetrics {
		c.updateWorkflowMetricEnd(name, duration, err)
	}
}

// 并发控制方法
func (c *StandardCollector) IncrementConcurrentActivities() {
	atomic.AddInt64(&c.metrics.ConcurrentActivities, 1)
}

func (c *StandardCollector) DecrementConcurrentActivities() {
	atomic.AddInt64(&c.metrics.ConcurrentActivities, -1)
}

func (c *StandardCollector) IncrementConcurrentWorkflows() {
	atomic.AddInt64(&c.metrics.ConcurrentWorkflows, 1)
}

func (c *StandardCollector) DecrementConcurrentWorkflows() {
	atomic.AddInt64(&c.metrics.ConcurrentWorkflows, -1)
}

// GetMetrics 获取指标（兼容旧接口，返回 map[string]int64）
func (c *StandardCollector) GetMetrics() map[string]int64 {
	// 获取最近1小时的数据（6个10分钟窗口）
	hourly := c.metrics.slidingWindow.getRecentMetrics(6)

	result := map[string]int64{
		// 小时数据
		"activities":            hourly["activities"],
		"workflows":             hourly["workflows"],
		"errors":                hourly["errors"],
		"activity_success":      hourly["activity_success"],
		"activity_failure":      hourly["activity_failure"],
		"workflow_success":      hourly["workflow_success"],
		"workflow_failure":      hourly["workflow_failure"],
		"activity_duration_ms":  hourly["activity_duration_ms"],
		"workflow_duration_ms":  hourly["workflow_duration_ms"],
		"max_activity_duration": hourly["max_activity_duration"],
		"max_workflow_duration": hourly["max_workflow_duration"],
		"activity_retries":      hourly["activity_retries"],
		"workflow_retries":      hourly["workflow_retries"],

		// 实时并发
		"concurrent_activities": atomic.LoadInt64(&c.metrics.ConcurrentActivities),
		"concurrent_workflows":  atomic.LoadInt64(&c.metrics.ConcurrentWorkflows),

		// 累积总数
		"total_activities": atomic.LoadInt64(&c.metrics.TotalActivityCount),
		"total_workflows":  atomic.LoadInt64(&c.metrics.TotalWorkflowCount),
		"total_errors":     atomic.LoadInt64(&c.metrics.TotalErrorCount),
	}

	// 计算成功率
	if total := result["activity_success"] + result["activity_failure"]; total > 0 {
		result["activity_success_rate"] = int64(float64(result["activity_success"]) / float64(total) * 100)
	}

	if total := result["workflow_success"] + result["workflow_failure"]; total > 0 {
		result["workflow_success_rate"] = int64(float64(result["workflow_success"]) / float64(total) * 100)
	}

	return result
}

// GetDetailedMetrics 获取详细指标
func (c *StandardCollector) GetDetailedMetrics() (map[string]*ExecutionMetric, map[string]*ExecutionMetric) {
	if !c.config.EnableDetailedMetrics {
		return nil, nil
	}

	activities := make(map[string]*ExecutionMetric)
	workflows := make(map[string]*ExecutionMetric)

	c.metrics.DetailedActivity.Range(func(key, value interface{}) bool {
		activities[key.(string)] = value.(*ExecutionMetric)
		return true
	})

	c.metrics.DetailedWorkflow.Range(func(key, value interface{}) bool {
		workflows[key.(string)] = value.(*ExecutionMetric)
		return true
	})

	return activities, workflows
}

// Reset 手动重置所有指标
func (c *StandardCollector) Reset() {
	c.resetAccumulatedMetrics()

	// 重置滑动窗口
	c.metrics.slidingWindow = NewSlidingWindowMetrics(
		c.config.MetricsWindowSize,
		c.config.MetricsWindowDuration,
	)

	// 重置并发计数
	atomic.StoreInt64(&c.metrics.ConcurrentActivities, 0)
	atomic.StoreInt64(&c.metrics.ConcurrentWorkflows, 0)
}

// 详细指标更新方法
func (c *StandardCollector) updateActivityMetricStart(name string, attempt int32) {
	val, _ := c.metrics.DetailedActivity.LoadOrStore(name, &ExecutionMetric{
		Name:        name,
		MinDuration: int64(^uint64(0) >> 1), // MaxInt64
	})

	metric := val.(*ExecutionMetric)
	if attempt > 1 {
		atomic.AddInt64(&metric.RetryCount, 1)
	}
}

func (c *StandardCollector) updateActivityMetricEnd(name string, duration int64, err error) {
	val, _ := c.metrics.DetailedActivity.LoadOrStore(name, &ExecutionMetric{
		Name:        name,
		MinDuration: int64(^uint64(0) >> 1), // MaxInt64
	})

	metric := val.(*ExecutionMetric)
	atomic.AddInt64(&metric.Count, 1)
	atomic.AddInt64(&metric.TotalDuration, duration)

	if err != nil {
		atomic.AddInt64(&metric.FailureCount, 1)
		metric.mu.Lock()
		metric.LastError = err.Error()
		metric.mu.Unlock()
	} else {
		atomic.AddInt64(&metric.SuccessCount, 1)
	}

	// 更新最大/最小执行时间
	c.updateMaxDuration(&metric.MaxDuration, duration)
	c.updateMinDuration(&metric.MinDuration, duration)

	// 更新最后执行时间
	metric.mu.Lock()
	metric.LastExecution = time.Now()
	metric.mu.Unlock()
}

func (c *StandardCollector) updateWorkflowMetricStart(name string) {
	_, _ = c.metrics.DetailedWorkflow.LoadOrStore(name, &ExecutionMetric{
		Name:        name,
		MinDuration: int64(^uint64(0) >> 1), // MaxInt64
	})
}

func (c *StandardCollector) updateWorkflowMetricEnd(name string, duration int64, err error) {
	val, _ := c.metrics.DetailedWorkflow.LoadOrStore(name, &ExecutionMetric{
		Name:        name,
		MinDuration: int64(^uint64(0) >> 1), // MaxInt64
	})

	metric := val.(*ExecutionMetric)
	atomic.AddInt64(&metric.Count, 1)
	atomic.AddInt64(&metric.TotalDuration, duration)

	if err != nil {
		atomic.AddInt64(&metric.FailureCount, 1)
		metric.mu.Lock()
		metric.LastError = err.Error()
		metric.mu.Unlock()
	} else {
		atomic.AddInt64(&metric.SuccessCount, 1)
	}

	// 更新最大/最小执行时间
	c.updateMaxDuration(&metric.MaxDuration, duration)
	c.updateMinDuration(&metric.MinDuration, duration)

	// 更新最后执行时间
	metric.mu.Lock()
	metric.LastExecution = time.Now()
	metric.mu.Unlock()
}

// 辅助方法
func (c *StandardCollector) updateMaxDuration(target *int64, duration int64) {
	for {
		current := atomic.LoadInt64(target)
		if duration <= current || atomic.CompareAndSwapInt64(target, current, duration) {
			break
		}
	}
}

func (c *StandardCollector) updateMinDuration(target *int64, duration int64) {
	for {
		current := atomic.LoadInt64(target)
		if duration >= current || atomic.CompareAndSwapInt64(target, current, duration) {
			break
		}
	}
}

// Stop 停止收集器
func (c *StandardCollector) Stop() {
	if c.resetTicker != nil {
		c.resetTicker.Stop()
	}
	select {
	case <-c.stopChan:
		// already closed
	default:
		close(c.stopChan)
	}
}

// ========================= 指标管理器 =========================

type Manager struct {
	collector Collector
	reporter  *Reporter
	config    *Config
}

var (
	globalManager     *Manager
	globalManagerOnce sync.Once
)

// GetGlobalManager 获取全局管理器
func GetGlobalManager() *Manager {
	globalManagerOnce.Do(func() {
		globalManager = NewManager(DefaultConfig())
	})
	return globalManager
}

// InitGlobalManager 初始化全局管理器
func InitGlobalManager(config *Config) {
	globalManagerOnce.Do(func() {
		globalManager = NewManager(config)
	})
}

// NewManager 创建管理器
func NewManager(config *Config) *Manager {
	if config == nil {
		config = DefaultConfig()
	}

	m := &Manager{
		collector: NewStandardCollector(config),
		config:    config,
	}

	// 启动自动上报
	if config.EnableAutoReport && config.ReportInterval > 0 {
		m.reporter = NewReporter(config, m.collector)
		m.reporter.Start()
	}

	return m
}

// GetCollector 获取收集器
func (m *Manager) GetCollector() Collector {
	return m.collector
}

// GetMetrics 获取指标
func (m *Manager) GetMetrics() map[string]int64 {
	return m.collector.GetMetrics()
}

// GetDetailedMetrics 获取详细指标
func (m *Manager) GetDetailedMetrics() (map[string]*ExecutionMetric, map[string]*ExecutionMetric) {
	return m.collector.GetDetailedMetrics()
}

// Reset 重置指标
func (m *Manager) Reset() {
	m.collector.Reset()
}

// Stop 停止管理器
func (m *Manager) Stop() {
	if m.reporter != nil {
		m.reporter.Stop()
	}
	if stopper, ok := m.collector.(*StandardCollector); ok {
		stopper.Stop()
	}
}

// ========================= 指标上报器 =========================

type Reporter struct {
	config    *Config
	collector Collector
	stopChan  chan struct{}
	wg        sync.WaitGroup
}

func NewReporter(config *Config, collector Collector) *Reporter {
	return &Reporter{
		config:    config,
		collector: collector,
		stopChan:  make(chan struct{}),
	}
}

func (r *Reporter) Start() {
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()

		ticker := time.NewTicker(r.config.ReportInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				r.report()
			case <-r.stopChan:
				return
			}
		}
	}()
}

func (r *Reporter) Stop() {
	select {
	case <-r.stopChan:
		// already closed
	default:
		close(r.stopChan)
	}
	r.wg.Wait()
}

func (r *Reporter) report() {
	metrics := r.collector.GetMetrics()

	// 自定义上报
	if r.config.CustomReporter != nil {
		// 转换为 map[string]interface{} 以支持更丰富的数据
		metricsInterface := make(map[string]interface{})
		for k, v := range metrics {
			metricsInterface[k] = v
		}
		r.config.CustomReporter(metricsInterface)
		return
	}

	// 默认日志上报
	if r.config.Logger != nil {
		r.config.Logger.Info("Metrics Report",
			"activities", metrics["activities"],
			"workflows", metrics["workflows"],
			"errors", metrics["errors"],
			"activity_success_rate", metrics["activity_success_rate"],
			"workflow_success_rate", metrics["workflow_success_rate"],
			"concurrent_activities", metrics["concurrent_activities"],
			"concurrent_workflows", metrics["concurrent_workflows"],
			"total_activities", metrics["total_activities"],
			"total_workflows", metrics["total_workflows"],
		)

		// 详细指标
		if r.config.EnableDetailedMetrics {
			r.reportDetailedMetrics()
		}
	}
}

func (r *Reporter) reportDetailedMetrics() {
	activities, workflows := r.collector.GetDetailedMetrics()

	// 报告活动指标
	for name, metric := range activities {
		if metric.Count == 0 {
			continue
		}

		avgDuration := metric.TotalDuration / metric.Count
		successRate := float64(metric.SuccessCount) / float64(metric.Count) * 100

		r.config.Logger.Debug("Activity Metric",
			"name", name,
			"count", metric.Count,
			"success_rate", successRate,
			"avg_duration_ms", avgDuration,
			"max_duration_ms", metric.MaxDuration,
			"min_duration_ms", metric.MinDuration,
			"retry_count", metric.RetryCount,
			"last_execution", metric.LastExecution.Format(time.RFC3339),
		)
	}

	// 报告工作流指标
	for name, metric := range workflows {
		if metric.Count == 0 {
			continue
		}

		avgDuration := metric.TotalDuration / metric.Count
		successRate := float64(metric.SuccessCount) / float64(metric.Count) * 100

		r.config.Logger.Debug("Workflow Metric",
			"name", name,
			"count", metric.Count,
			"success_rate", successRate,
			"avg_duration_ms", avgDuration,
			"max_duration_ms", metric.MaxDuration,
			"min_duration_ms", metric.MinDuration,
			"last_execution", metric.LastExecution.Format(time.RFC3339),
		)
	}
}

// ========================= 便捷函数 =========================

// GetGlobalCollector 获取全局收集器
func GetGlobalCollector() Collector {
	return GetGlobalManager().GetCollector()
}

// GetGlobalMetrics 获取全局指标
func GetGlobalMetrics() map[string]int64 {
	return GetGlobalManager().GetMetrics()
}

// ResetGlobalMetrics 重置全局指标
func ResetGlobalMetrics() {
	GetGlobalManager().Reset()
}

// StopGlobalManager 停止全局管理器
func StopGlobalManager() {
	GetGlobalManager().Stop()
}
