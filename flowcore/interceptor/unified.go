package interceptor

import (
	"context"
	"fmt"
	"time"

	"github.com/hanyougame/glib/flowcore/logger"
	flowcorepkg "github.com/hanyougame/glib/flowcore/pkg"
	"github.com/hanyougame/glib/metadata"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
	commonpb "go.temporal.io/api/common/v1"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"
)

// ========================= 配置 =========================

// Config 拦截器配置
type Config struct {
	EnableTrace      bool              // 启用链路追踪
	EnableMetrics    bool              // 启用指标收集
	Logger           log.Logger        // 日志器
	CustomAttributes map[string]string // 自定义属性
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		EnableTrace:      true,
		EnableMetrics:    true,
		CustomAttributes: make(map[string]string),
	}
}

// ========================= 统一拦截器 =========================

// UnifiedInterceptor 统一的拦截器，整合 trace、metrics 和监控功能
type UnifiedInterceptor struct {
	interceptor.InterceptorBase
	config    *Config
	collector flowcorepkg.Collector
}

// NewUnifiedInterceptor 创建统一拦截器
func NewUnifiedInterceptor(cfg *Config) interceptor.Interceptor {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	if cfg.Logger == nil {
		cfg.Logger = logger.NewDevelopmentLogger("info")
	}

	ui := &UnifiedInterceptor{
		config: cfg,
	}

	// 如果启用指标，获取全局收集器
	if cfg.EnableMetrics {
		ui.collector = flowcorepkg.GetGlobalCollector()
	}

	return ui
}

// InterceptActivity 拦截 Activity 调用
func (u *UnifiedInterceptor) InterceptActivity(
	ctx context.Context,
	next interceptor.ActivityInboundInterceptor,
) interceptor.ActivityInboundInterceptor {
	i := &activityInbound{
		config:    u.config,
		collector: u.collector,
	}
	i.Next = next
	return i
}

// InterceptWorkflow 拦截 Workflow 调用
func (u *UnifiedInterceptor) InterceptWorkflow(
	ctx workflow.Context,
	next interceptor.WorkflowInboundInterceptor,
) interceptor.WorkflowInboundInterceptor {
	info := workflow.GetInfo(ctx)
	_logger := workflow.GetLogger(ctx)
	// 如果 logger 为 nil，使用配置中的默认 logger 作为后备
	if _logger == nil && u.config != nil && u.config.Logger != nil {
		_logger = u.config.Logger
	}

	i := &workflowInbound{
		config:    u.config,
		collector: u.collector,
		info:      info,
		logger:    _logger,
	}
	i.Next = next
	return i
}

// GetMetrics 获取当前指标
func (u *UnifiedInterceptor) GetMetrics() map[string]int64 {
	if u.collector != nil {
		return u.collector.GetMetrics()
	}
	return nil
}

// GetDetailedMetrics 获取详细指标
func (u *UnifiedInterceptor) GetDetailedMetrics() (map[string]*flowcorepkg.ExecutionMetric, map[string]*flowcorepkg.ExecutionMetric) {
	if u.collector != nil {
		return u.collector.GetDetailedMetrics()
	}
	return nil, nil
}

// ========================= Activity 拦截器 =========================

type activityInbound struct {
	interceptor.ActivityInboundInterceptorBase
	config    *Config
	collector flowcorepkg.Collector
}

func (a *activityInbound) Init(outbound interceptor.ActivityOutboundInterceptor) error {
	i := &activityOutbound{
		config: a.config,
	}
	i.Next = outbound
	return a.Next.Init(i)
}

func (a *activityInbound) ExecuteActivity(
	ctx context.Context,
	in *interceptor.ExecuteActivityInput,
) (interface{}, error) {
	// 获取 Activity 信息
	info := activity.GetInfo(ctx)
	activityName := info.ActivityType.Name

	a.config.Logger.Info("[Activity Interceptor] Starting activity execution",
		"activity", activityName,
		"workflow_id", info.WorkflowExecution.ID,
		"run_id", info.WorkflowExecution.RunID,
		"activity_id", info.ActivityID,
		"attempt", info.Attempt,
		"task_queue", info.TaskQueue,
	)

	// 增加并发计数
	if a.collector != nil {
		a.collector.IncrementConcurrentActivities()
		defer a.collector.DecrementConcurrentActivities()
	}

	// Metrics: 记录开始时间
	var start time.Time
	if a.config.EnableMetrics && a.collector != nil {
		start = time.Now()
		a.collector.RecordActivityStart(activityName, info.Attempt)

		a.config.Logger.Debug("Activity started",
			"activity", activityName,
			"workflow_id", info.WorkflowExecution.ID,
			"activity_id", info.ActivityID,
			"attempt", info.Attempt,
		)
	}

	// Trace: 从 Header 恢复 trace context
	var span oteltrace.Span
	if a.config.EnableTrace {
		// 记录恢复前的状态
		a.config.Logger.Info("[Activity Interceptor] Restoring context from header", "activity", activityName)

		ctx = a.restoreContextFromHeader(ctx)
		ctx, span = a.createActivitySpan(ctx, info)

		if span != nil {
			defer func() {
				if r := recover(); r != nil {
					span.RecordError(fmt.Errorf("panic: %v", r))
					span.SetStatus(codes.Error, fmt.Sprintf("panic: %v", r))
					span.End()
					panic(r)
				}
				span.End()
			}()
		}
	}

	// 执行 Activity
	result, err := a.Next.ExecuteActivity(ctx, in)

	// 记录执行时间
	duration := time.Since(start).Milliseconds()

	// Trace: 记录错误
	if span != nil {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, "success")
		}
		span.SetAttributes(attribute.Int64("duration_ms", duration))
	}

	// Metrics: 记录结束和持续时间
	if a.config.EnableMetrics && a.collector != nil {
		a.collector.RecordActivityEnd(activityName, duration, err)

		if err != nil {
			a.config.Logger.Error("Activity failed",
				"activity", activityName,
				"error", err,
				"duration_ms", duration,
				"attempt", info.Attempt,
			)
		} else {
			a.config.Logger.Debug("Activity completed",
				"activity", activityName,
				"duration_ms", duration,
			)
		}
	}

	return result, err
}

// restoreContextFromHeader 从 Header 恢复上下文
func (a *activityInbound) restoreContextFromHeader(ctx context.Context) context.Context {
	header := interceptor.Header(ctx)
	if header == nil {
		a.config.Logger.Warn("[Activity Interceptor] No header found in context")
		return ctx
	}

	contextData := extractContextFromHeader(header)
	if contextData == nil {
		a.config.Logger.Warn("[Activity Interceptor] No context data found in header")
		return ctx
	}

	// 记录恢复的上下文信息
	a.config.Logger.Info("[Activity Interceptor] Context data extracted from header",
		"trace_id", contextData.TraceID,
		"span_id", contextData.SpanID,
		"merchant_id", contextData.MerchantID,
		"currency_code", contextData.CurrencyCode,
		"user_id", contextData.UserID,
		"merchant_user_id", contextData.MerchantUserID,
		"baggage_count", len(contextData.Baggage),
	)

	//// 恢复 trace context
	//if contextData.TraceID != "" {
	//	ctx = tracex.RestoreTraceContext(ctx, contextData.TraceID, contextData.SpanID)
	//}

	// 恢复商户信息
	if contextData.MerchantID > 0 {
		//ctx = metadata.WithMetadata(ctx, metadata.CtxMerchantID, contextData.MerchantID)
		ctx = metadata.WithMetadata(ctx, metadata.CtxCurrencyCode, contextData.CurrencyCode)
	}

	if contextData.UserID > 0 {
		ctx = metadata.WithMetadata(ctx, metadata.CtxJWTUserId, contextData.UserID)
	}

	return ctx
}

// createActivitySpan 创建 Activity 的 trace span
func (a *activityInbound) createActivitySpan(ctx context.Context, info activity.Info) (context.Context, oteltrace.Span) {
	tracer := otel.Tracer("temporal-activity")

	// 添加自定义属性
	attributes := []attribute.KeyValue{
		attribute.String("activity.type", info.ActivityType.Name),
		attribute.String("activity.id", info.ActivityID),
		attribute.String("workflow.id", info.WorkflowExecution.ID),
		attribute.String("workflow.run_id", info.WorkflowExecution.RunID),
		attribute.String("task_queue", info.TaskQueue),
		attribute.Int("attempt", int(info.Attempt)),
	}

	// 添加自定义属性
	for k, v := range a.config.CustomAttributes {
		attributes = append(attributes, attribute.String(k, v))
	}

	ctx, span := tracer.Start(ctx,
		fmt.Sprintf("activity.%s", info.ActivityType.Name),
		oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		oteltrace.WithAttributes(attributes...))

	if span != nil {
		spanCtx := span.SpanContext()
		a.config.Logger.Info("[Activity Interceptor] Span created",
			"activity", info.ActivityType.Name,
			"trace_id", spanCtx.TraceID().String(),
			"span_id", spanCtx.SpanID().String(),
		)
	}

	return ctx, span
}

type activityOutbound struct {
	interceptor.ActivityOutboundInterceptorBase
	config *Config
}

func (a *activityOutbound) GetLogger(ctx context.Context) log.Logger {
	_logger := a.Next.GetLogger(ctx)

	// 如果有 trace span，增强日志
	if a.config.EnableTrace {
		if span := oteltrace.SpanFromContext(ctx); span != nil && span.SpanContext().IsValid() {
			spanCtx := span.SpanContext()
			return log.With(_logger,
				"trace", spanCtx.TraceID().String(),
				"span", spanCtx.SpanID().String(),
			)
		}
	}
	return _logger
}

// ========================= Workflow 拦截器 =========================

type workflowInbound struct {
	interceptor.WorkflowInboundInterceptorBase
	config    *Config
	collector flowcorepkg.Collector
	info      *workflow.Info
	logger    log.Logger
}

func (w *workflowInbound) Init(outbound interceptor.WorkflowOutboundInterceptor) error {
	// 确保 logger 可用
	if w.logger == nil && w.config != nil && w.config.Logger != nil {
		w.logger = w.config.Logger
	}

	i := &workflowOutbound{
		config: w.config,
		logger: w.logger,
	}
	i.Next = outbound
	return w.Next.Init(i)
}

func (w *workflowInbound) ExecuteWorkflow(
	ctx workflow.Context,
	in *interceptor.ExecuteWorkflowInput,
) (interface{}, error) {
	// 延迟获取 workflow info，如果之前没有获取到
	if w.info == nil {
		w.info = workflow.GetInfo(ctx)
	}

	// 确保 logger 可用
	if w.logger == nil {
		w.logger = workflow.GetLogger(ctx)
		if w.logger == nil && w.config != nil && w.config.Logger != nil {
			w.logger = w.config.Logger
		}
	}

	// 安全地获取 workflow 名称和相关信息
	workflowName, workflowID, runID, parentWorkflowID := w.extractWorkflowInfo()

	// 增加并发计数
	if w.collector != nil {
		w.collector.IncrementConcurrentWorkflows()
		defer w.collector.DecrementConcurrentWorkflows()
	}

	// Metrics: 记录开始时间
	var start time.Time
	if w.config.EnableMetrics && w.collector != nil {
		start = workflow.Now(ctx)
		w.collector.RecordWorkflowStart(workflowName)
		w.logWorkflowStart(workflowName, workflowID, runID, parentWorkflowID)
	}

	// Trace: 从 Memo 中提取 trace 信息并存储到 workflow context
	if w.config.EnableTrace {
		ctx = w.setupWorkflowTrace(ctx, workflowName, workflowID)
	}

	// 执行 Workflow
	result, err := w.Next.ExecuteWorkflow(ctx, in)

	// 记录执行时间
	duration := workflow.Now(ctx).Sub(start).Milliseconds()

	// Metrics: 记录结束和持续时间
	if w.config.EnableMetrics && w.collector != nil {
		w.collector.RecordWorkflowEnd(workflowName, duration, err)
		w.logWorkflowEnd(workflowName, duration, err)
	}

	return result, err
}

// extractWorkflowInfo 安全地提取 workflow 信息
func (w *workflowInbound) extractWorkflowInfo() (name, id, runID, parentID string) {
	name = "unknown"

	if w.info != nil {
		if w.info.WorkflowType.Name != "" {
			name = w.info.WorkflowType.Name
		}
		id = w.info.WorkflowExecution.ID
		runID = w.info.WorkflowExecution.RunID

		if w.info.ParentWorkflowExecution != nil {
			parentID = w.info.ParentWorkflowExecution.ID
		}
	}

	return
}

// setupWorkflowTrace 设置 workflow 的 trace 上下文
func (w *workflowInbound) setupWorkflowTrace(ctx workflow.Context, workflowName, workflowID string) workflow.Context {
	// 需要多重空值检查
	if w.info == nil || w.info.Memo == nil || w.info.Memo.Fields == nil {
		if w.logger != nil {
			w.logger.Warn("[Workflow Interceptor] No Memo found in workflow info",
				"has_info", w.info != nil,
				"has_memo", w.info != nil && w.info.Memo != nil,
				"has_fields", w.info != nil && w.info.Memo != nil && w.info.Memo.Fields != nil,
			)
		}
		return ctx
	}

	// 记录 Memo 字段
	if w.logger != nil {
		fieldKeys := make([]string, 0, len(w.info.Memo.Fields))
		for k := range w.info.Memo.Fields {
			fieldKeys = append(fieldKeys, k)
		}
		w.logger.Info("[Workflow Interceptor] Memo fields found",
			"field_count", len(w.info.Memo.Fields),
			"field_keys", fieldKeys,
		)
	}

	contextData := flowcorepkg.ExtractContextFromMemo(w.info.Memo)

	if w.logger != nil {
		if contextData != nil {
			w.logger.Info("[Workflow Interceptor] Context data extracted from Memo",
				"trace_id", contextData.TraceID,
				"span_id", contextData.SpanID,
				"merchant_id", contextData.MerchantID,
				"currency_code", contextData.CurrencyCode,
				"user_id", contextData.UserID,
				"merchant_user_id", contextData.MerchantUserID,
				"baggage", contextData.Baggage,
				"is_valid", contextData.IsValid(),
			)
		}
	}

	// 检查 contextData 不为 nil 并且有效
	if contextData != nil && contextData.IsValid() {
		ctx = flowcorepkg.WorkflowContextWithData(ctx, contextData)

		if w.logger != nil {
			w.logger.Info("[Workflow Interceptor] Context data stored in workflow context",
				"workflow_type", workflowName,
				"workflow_id", workflowID,
				"trace_id", contextData.TraceID,
				"span_id", contextData.SpanID,
				"merchant_id", contextData.MerchantID,
				"currency_code", contextData.CurrencyCode,
			)
		}
		// 验证存储是否成功  立即尝试读取以验证
		verifyData := flowcorepkg.DataFromWorkflowContext(ctx)
		if verifyData != nil {
			w.logger.Info("[Workflow Interceptor] Context data verification successful",
				"stored_trace_id", verifyData.TraceID,
				"stored_merchant_id", verifyData.MerchantID,
			)
		}
	}

	return ctx
}

// logWorkflowStart 记录 workflow 开始日志
func (w *workflowInbound) logWorkflowStart(name, id, runID, parentID string) {
	if w.logger == nil {
		return
	}

	// 构建日志字段，只添加非空值
	logFields := []interface{}{"workflow", name}
	if id != "" {
		logFields = append(logFields, "workflow_id", id)
	}
	if runID != "" {
		logFields = append(logFields, "run_id", runID)
	}
	if parentID != "" {
		logFields = append(logFields, "parent_workflow_id", parentID)
	}

	w.logger.Debug("Workflow started", logFields...)
}

// logWorkflowEnd 记录 workflow 结束日志
func (w *workflowInbound) logWorkflowEnd(name string, duration int64, err error) {
	if w.logger == nil {
		return
	}

	if err != nil {
		w.logger.Error("Workflow failed",
			"workflow", name,
			"error", err,
			"duration_ms", duration,
		)
	} else {
		w.logger.Debug("Workflow completed",
			"workflow", name,
			"duration_ms", duration,
		)
	}
}

type workflowOutbound struct {
	interceptor.WorkflowOutboundInterceptorBase
	config *Config
	logger log.Logger
}

func (w *workflowOutbound) ExecuteActivity(
	ctx workflow.Context,
	activityType string,
	args ...interface{},
) workflow.Future {
	w.injectContextToHeader(ctx, "activity")
	return w.Next.ExecuteActivity(ctx, activityType, args...)
}

func (w *workflowOutbound) ExecuteLocalActivity(
	ctx workflow.Context,
	activityType string,
	args ...interface{},
) workflow.Future {
	w.injectContextToHeader(ctx, "local activity")
	return w.Next.ExecuteLocalActivity(ctx, activityType, args...)
}

func (w *workflowOutbound) ExecuteChildWorkflow(
	ctx workflow.Context,
	childWorkflowType string,
	args ...interface{},
) workflow.ChildWorkflowFuture {
	w.injectContextToHeader(ctx, "child workflow")
	return w.Next.ExecuteChildWorkflow(ctx, childWorkflowType, args...)
}

// injectContextToHeader 将上下文注入到 header
func (w *workflowOutbound) injectContextToHeader(ctx workflow.Context, executeType string) {
	if !w.config.EnableTrace {
		return
	}

	header := interceptor.WorkflowHeader(ctx)
	if header == nil {
		return
	}

	contextData := flowcorepkg.DataFromWorkflowContext(ctx)

	if w.logger != nil {
		if contextData != nil {
			w.logger.Info("[Workflow Outbound] Injecting context to header",
				"execute_type", executeType,

				"trace_id", contextData.TraceID,
				"span_id", contextData.SpanID,
				"merchant_id", contextData.MerchantID,
				"user_id", contextData.UserID,
				"baggage_count", len(contextData.Baggage),
			)
		}
	}

	if err := writeContextToHeader(contextData, header); err != nil {
		// 安全地获取和使用 logger
		_logger := workflow.GetLogger(ctx)
		if _logger == nil {
			_logger = w.logger
		}
		if _logger != nil {
			_logger.Error("Failed to write context to header",
				"type", executeType,
				"error", err)
		}
	}
}

func (w *workflowOutbound) GetLogger(ctx workflow.Context) log.Logger {
	loggerVal := w.Next.GetLogger(ctx)

	// 如果从 Next 获取的 logger 为 nil，使用备用 logger
	if loggerVal == nil {
		if w.logger != nil {
			loggerVal = w.logger
		} else if w.config != nil && w.config.Logger != nil {
			loggerVal = w.config.Logger
		}
	}

	// 如果还是没有 logger，直接返回 nil（让上层处理）
	if loggerVal == nil {
		return nil
	}

	// 如果启用了 trace，增强日志
	if w.config.EnableTrace {
		loggerVal = w.enhanceLoggerWithTrace(ctx, loggerVal)
	}

	return loggerVal
}

// enhanceLoggerWithTrace 使用 trace 信息增强 logger
func (w *workflowOutbound) enhanceLoggerWithTrace(ctx workflow.Context, logger log.Logger) log.Logger {
	fields := make(map[string]interface{})

	// 使用新的 API 获取 trace 信息
	traceID, spanID := flowcorepkg.GetTraceFromWorkflow(ctx)
	if traceID != "" {
		fields["trace"] = traceID
	}
	if spanID != "" {
		fields["span"] = spanID
	}

	// 添加商户信息
	merchantID, currencyCode := flowcorepkg.GetMerchantFromWorkflow(ctx)
	if merchantID > 0 {
		fields["merchant_id"] = merchantID
		fields["currency_code"] = currencyCode
	}

	// 添加自定义属性到日志
	for k, v := range w.config.CustomAttributes {
		fields[k] = v
	}

	if len(fields) > 0 {
		for k, v := range fields {
			logger = log.With(logger, k, v)
		}
	}

	return logger
}

// ========================= 辅助函数 =========================

const headerKey = "unified-context"

// extractContextFromHeader 从 Header 中提取上下文
func extractContextFromHeader(header map[string]*commonpb.Payload) *flowcorepkg.ContextData {
	payload := header[headerKey]
	if payload == nil {
		return nil
	}

	var contextData flowcorepkg.ContextData
	if err := converter.GetDefaultDataConverter().FromPayload(payload, &contextData); err != nil {
		return nil
	}

	return &contextData
}

// writeContextToHeader 将上下文写入 Header
func writeContextToHeader(contextData *flowcorepkg.ContextData, header map[string]*commonpb.Payload) error {
	if contextData == nil {
		return nil
	}

	payload, err := converter.GetDefaultDataConverter().ToPayload(contextData)
	if err != nil {
		return err
	}

	header[headerKey] = payload
	return nil
}
