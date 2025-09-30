package pkg

import (
	"context"
	"time"

	"github.com/hanyougame/glib/metadata"
	"github.com/spf13/cast"
	commonpb "go.temporal.io/api/common/v1"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/workflow"
)

// ========================= Context Keys =========================

// contextKey 用于在 context 中存储追踪信息的键类型
type contextKey string

const (
	contextKeyTraceID        contextKey = "trace"
	contextKeySpanID         contextKey = "span"
	contextKeyMerchantID     contextKey = "merchant_id"
	contextKeyCurrencyCode   contextKey = "currency_code"
	contextKeyMerchantUserID contextKey = "merchant_user_id"
	contextKeyUserID         contextKey = "user_id"
	contextKeyBaggage        contextKey = "baggage" // 附加信息
)

// ========================= Context Data =========================

// ContextData 在 Workflow 和 Activity 之间传递的完整上下文数据
type ContextData struct {
	// 追踪信息
	TraceID string `json:"trace_id,omitempty"`
	SpanID  string `json:"span_id,omitempty"`

	// 业务信息
	MerchantID     int64  `json:"merchant_id,omitempty"`
	CurrencyCode   string `json:"currency_code,omitempty"`
	MerchantUserID string `json:"merchant_user_id,omitempty"`
	UserID         int64  `json:"user_id,omitempty"`

	// 附加数据
	Baggage map[string]string `json:"baggage,omitempty"`

	// 时间信息
	Timestamp time.Time `json:"timestamp,omitempty"`
}

// NewContextData 创建新的上下文数据
func NewContextData() *ContextData {
	return &ContextData{
		Baggage:   make(map[string]string),
		Timestamp: time.Now(),
	}
}

// WithTraceInfo 设置追踪信息
func (c *ContextData) WithTraceInfo(traceID, spanID string) *ContextData {
	c.TraceID = traceID
	c.SpanID = spanID
	return c
}

// WithMerchantInfo 设置商户信息
func (c *ContextData) WithMerchantInfo(merchantID int64, currencyCode string) *ContextData {
	c.MerchantID = merchantID
	c.CurrencyCode = currencyCode
	return c
}

// WithUserInfo 设置用户信息
func (c *ContextData) WithUserInfo(userID int64, merchantUserID string) *ContextData {
	c.UserID = userID
	c.MerchantUserID = merchantUserID
	return c
}

// AddBaggage 添加 Baggage 数据
func (c *ContextData) AddBaggage(key, value string) *ContextData {
	if c.Baggage == nil {
		c.Baggage = make(map[string]string)
	}
	c.Baggage[key] = value
	return c
}

// GetBaggage 获取 Baggage 数据
func (c *ContextData) GetBaggage(key string) (string, bool) {
	if c.Baggage == nil {
		return "", false
	}
	val, ok := c.Baggage[key]
	return val, ok
}

// Clone 克隆上下文数据
func (c *ContextData) Clone() *ContextData {
	if c == nil {
		return NewContextData()
	}

	cloned := &ContextData{
		TraceID:        c.TraceID,
		SpanID:         c.SpanID,
		MerchantID:     c.MerchantID,
		CurrencyCode:   c.CurrencyCode,
		MerchantUserID: c.MerchantUserID,
		UserID:         c.UserID,
		Timestamp:      c.Timestamp,
	}

	if c.Baggage != nil {
		cloned.Baggage = make(map[string]string)
		for k, v := range c.Baggage {
			cloned.Baggage[k] = v
		}
	}

	return cloned
}

// IsValid 检查上下文数据是否有效
func (c *ContextData) IsValid() bool {
	return c != nil && c.TraceID != ""
}

// ========================= Standard Context Operations =========================

// ContextWithData 将 ContextData 存储到 context.Context 中
func ContextWithData(ctx context.Context, data *ContextData) context.Context {
	if data == nil {
		return ctx
	}

	// 存储追踪信息
	if data.TraceID != "" {
		ctx = context.WithValue(ctx, contextKeyTraceID, data.TraceID)
	}
	if data.SpanID != "" {
		ctx = context.WithValue(ctx, contextKeySpanID, data.SpanID)
	}

	// 存储业务信息
	if data.MerchantID > 0 {
		ctx = context.WithValue(ctx, contextKeyMerchantID, data.MerchantID)
	}
	if data.CurrencyCode != "" {
		ctx = context.WithValue(ctx, contextKeyCurrencyCode, data.CurrencyCode)
	}
	if data.MerchantUserID != "" {
		ctx = context.WithValue(ctx, contextKeyMerchantUserID, data.MerchantUserID)
	}
	if data.UserID > 0 {
		ctx = context.WithValue(ctx, contextKeyUserID, data.UserID)
	}

	// 存储附加信息
	if len(data.Baggage) > 0 {
		ctx = context.WithValue(ctx, contextKeyBaggage, data.Baggage)
	}

	return ctx
}

// DataFromContext 从 context.Context 中提取 ContextData
func DataFromContext(ctx context.Context) *ContextData {
	data := NewContextData()

	// 提取追踪信息
	if val := ctx.Value(contextKeyTraceID); val != nil {
		data.TraceID, _ = val.(string)
	}
	if val := ctx.Value(contextKeySpanID); val != nil {
		data.SpanID, _ = val.(string)
	}

	// 提取业务信息
	if val := ctx.Value(contextKeyMerchantID); val != nil {
		data.MerchantID, _ = val.(int64)
	}
	if val := ctx.Value(contextKeyCurrencyCode); val != nil {
		data.CurrencyCode, _ = val.(string)
	}
	if val := ctx.Value(contextKeyMerchantUserID); val != nil {
		data.MerchantUserID, _ = val.(string)
	}
	if val := ctx.Value(contextKeyUserID); val != nil {
		data.UserID, _ = val.(int64)
	}

	// 提取附加信息
	if val := ctx.Value(contextKeyBaggage); val != nil {
		data.Baggage, _ = val.(map[string]string)
	}

	return data
}

// ========================= Workflow Context Operations =========================

// WorkflowContextWithData 将 ContextData 存储到 workflow.Context 中
func WorkflowContextWithData(ctx workflow.Context, data *ContextData) workflow.Context {
	if data == nil {
		return ctx
	}

	// 存储追踪信息
	if data.TraceID != "" {
		ctx = workflow.WithValue(ctx, string(contextKeyTraceID), data.TraceID)
	}
	if data.SpanID != "" {
		ctx = workflow.WithValue(ctx, string(contextKeySpanID), data.SpanID)
	}

	// 存储业务信息
	if data.MerchantID > 0 {
		ctx = workflow.WithValue(ctx, string(contextKeyMerchantID), data.MerchantID)
	}
	if data.CurrencyCode != "" {
		ctx = workflow.WithValue(ctx, string(contextKeyCurrencyCode), data.CurrencyCode)
	}
	if data.MerchantUserID != "" {
		ctx = workflow.WithValue(ctx, string(contextKeyMerchantUserID), data.MerchantUserID)
	}
	if data.UserID > 0 {
		ctx = workflow.WithValue(ctx, string(contextKeyUserID), data.UserID)
	}

	// 存储附加信息
	if len(data.Baggage) > 0 {
		ctx = workflow.WithValue(ctx, string(contextKeyBaggage), data.Baggage)
	}

	return ctx
}

// DataFromWorkflowContext 从 workflow.Context 中提取 ContextData
func DataFromWorkflowContext(ctx workflow.Context) *ContextData {
	data := NewContextData()

	// 提取追踪信息
	if val := ctx.Value(string(contextKeyTraceID)); val != nil {
		data.TraceID, _ = val.(string)
	}
	if val := ctx.Value(string(contextKeySpanID)); val != nil {
		data.SpanID, _ = val.(string)
	}

	// 提取业务信息
	if val := ctx.Value(string(contextKeyMerchantID)); val != nil {
		data.MerchantID, _ = val.(int64)
	}
	if val := ctx.Value(string(contextKeyCurrencyCode)); val != nil {
		data.CurrencyCode, _ = val.(string)
	}
	if val := ctx.Value(string(contextKeyMerchantUserID)); val != nil {
		data.MerchantUserID, _ = val.(string)
	}
	if val := ctx.Value(string(contextKeyUserID)); val != nil {
		data.UserID, _ = val.(int64)
	}

	// 提取附加信息
	if val := ctx.Value(string(contextKeyBaggage)); val != nil {
		data.Baggage, _ = val.(map[string]string)
	}

	// 如果 context 中没有，尝试从 Memo 获取
	if data.TraceID == "" {
		info := workflow.GetInfo(ctx)
		if info.Memo != nil {
			memoData := ExtractContextFromMemo(info.Memo)
			if data.TraceID == "" {
				data = memoData
			}
		}
	}

	return data
}

// ========================= Memo Operations =========================

// ExtractContextFromMemo 从 Memo 中提取上下文信息
func ExtractContextFromMemo(memo *commonpb.Memo) *ContextData {
	if memo == nil || memo.Fields == nil {
		return NewContextData()
	}

	dc := converter.GetDefaultDataConverter()
	contextData := NewContextData()

	if payload, ok := memo.Fields[metadata.CtxCurrencyCode]; ok {
		var currencyCode string
		if err := dc.FromPayload(payload, &currencyCode); err == nil {
			contextData.CurrencyCode = currencyCode
		}
	}

	if payload, ok := memo.Fields[metadata.CtxJWTUserId]; ok {
		var userIDStr string
		if err := dc.FromPayload(payload, &userIDStr); err == nil {
			contextData.UserID, _ = cast.ToInt64E(userIDStr)
		}
	}

	return contextData
}

// BuildMemoMapFromContextData 从 ContextData 构建 Memo 的字符串映射
// 返回一个 map[string]string，用于构建 Memo
func BuildMemoMapFromContextData(data *ContextData) map[string]string {
	if data == nil {
		return make(map[string]string)
	}

	memoMap := make(map[string]string)

	if data.CurrencyCode != "" {
		memoMap[metadata.CtxCurrencyCode] = data.CurrencyCode
	}

	// 添加用户信息
	if data.UserID > 0 {
		memoMap[metadata.CtxJWTUserId] = cast.ToString(data.UserID)
	}

	// 添加时间戳
	if !data.Timestamp.IsZero() {
		memoMap["timestamp"] = data.Timestamp.Format(time.RFC3339)
	}

	return memoMap
}

// BuildMemoFromContext 从 ContextData 构建 Memo
func BuildMemoFromContext(data *ContextData) (*commonpb.Memo, error) {
	if data == nil {
		return nil, nil
	}

	dc := converter.GetDefaultDataConverter()
	memo := &commonpb.Memo{
		Fields: make(map[string]*commonpb.Payload),
	}

	if data.CurrencyCode != "" {
		payload, err := dc.ToPayload(data.CurrencyCode)
		if err != nil {
			return nil, err
		}
		memo.Fields[metadata.CtxCurrencyCode] = payload
	}

	if data.UserID > 0 {
		payload, err := dc.ToPayload(cast.ToString(data.UserID))
		if err != nil {
			return nil, err
		}
		memo.Fields[metadata.CtxJWTUserId] = payload
	}

	return memo, nil
}

// MergeMemos 合并两个 Memo
func MergeMemos(memo1, memo2 *commonpb.Memo) *commonpb.Memo {
	if memo1 == nil {
		return memo2
	}
	if memo2 == nil {
		return memo1
	}

	merged := &commonpb.Memo{
		Fields: make(map[string]*commonpb.Payload),
	}

	// 复制 memo1 的字段
	for k, v := range memo1.Fields {
		merged.Fields[k] = v
	}

	// 复制 memo2 的字段（会覆盖 memo1 的同名字段）
	for k, v := range memo2.Fields {
		merged.Fields[k] = v
	}

	return merged
}

// ========================= Workflow Helper Functions =========================

// GetTraceFromWorkflow 从 Workflow 中获取 trace 信息
func GetTraceFromWorkflow(ctx workflow.Context) (traceID, spanID string) {
	data := DataFromWorkflowContext(ctx)
	return data.TraceID, data.SpanID
}

// GetMerchantFromWorkflow 从 Workflow 中获取商户信息
func GetMerchantFromWorkflow(ctx workflow.Context) (merchantID int64, currencyCode string) {
	data := DataFromWorkflowContext(ctx)
	return data.MerchantID, data.CurrencyCode
}

// GetUserFromWorkflow 从 Workflow 中获取用户信息
func GetUserFromWorkflow(ctx workflow.Context) (userID int64, merchantUserID string) {
	data := DataFromWorkflowContext(ctx)
	return data.UserID, data.MerchantUserID
}

// GetBaggageFromWorkflow 从 Workflow 中获取 Baggage
func GetBaggageFromWorkflow(ctx workflow.Context) map[string]string {
	data := DataFromWorkflowContext(ctx)
	return data.Baggage
}

// SetBaggageInWorkflow 在 Workflow 中设置 Baggage
func SetBaggageInWorkflow(ctx workflow.Context, key, value string) workflow.Context {
	data := DataFromWorkflowContext(ctx)
	data.AddBaggage(key, value)
	return WorkflowContextWithData(ctx, data)
}

// ========================= Activity Helper Functions =========================

// GetTraceFromActivity 从 Activity Context 中获取 trace 信息
func GetTraceFromActivity(ctx context.Context) (traceID, spanID string) {
	data := DataFromContext(ctx)
	return data.TraceID, data.SpanID
}

// GetMerchantFromActivity 从 Activity Context 中获取商户信息
func GetMerchantFromActivity(ctx context.Context) (merchantID int64, currencyCode string) {
	data := DataFromContext(ctx)
	return data.MerchantID, data.CurrencyCode
}

// GetUserFromActivity 从 Activity Context 中获取用户信息
func GetUserFromActivity(ctx context.Context) (userID int64, merchantUserID string) {
	data := DataFromContext(ctx)
	return data.UserID, data.MerchantUserID
}

// GetBaggageFromActivity 从 Activity Context 中获取 Baggage
func GetBaggageFromActivity(ctx context.Context) map[string]string {
	data := DataFromContext(ctx)
	return data.Baggage
}

// SetBaggageInActivity 在 Activity Context 中设置 Baggage
func SetBaggageInActivity(ctx context.Context, key, value string) context.Context {
	data := DataFromContext(ctx)
	data.AddBaggage(key, value)
	return ContextWithData(ctx, data)
}
