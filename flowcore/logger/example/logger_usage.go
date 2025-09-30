package example

import (
	"context"
	"time"

	"github.com/hanyougame/glib/flowcore/logger"
	"go.temporal.io/sdk/workflow"
)

// ========================= 基础使用 =========================

func BasicUsageExample() {
	// 1. 创建基础 logger
	log := logger.NewLogger("debug")
	log.Info("应用启动", "version", "1.0.0", "env", "production")

	// 2. 创建开发环境 logger（带颜色和图标）
	devLog := logger.NewDevelopmentLogger("debug")
	devLog.Debug("调试信息", "user_id", 12345)
	devLog.Info("用户登录成功", "username", "alice")
	devLog.Warn("请求延迟较高", "latency_ms", 500)
	devLog.Error("数据库连接失败", "error", "connection timeout")

	// 3. 创建生产环境 logger（JSON 格式）
	prodLog := logger.NewProductionLogger()
	prodLog.Info("订单创建", "order_id", "ORD-123456", "amount", 99.99)

	// 4. 使用 WithFields 添加固定字段
	userLog := log.WithFields(map[string]interface{}{
		"service": "user-service",
		"region":  "us-east-1",
	})
	userLog.Info("处理用户请求", "action", "get_profile")
}

// ========================= Workflow 中使用 =========================

// OrderWorkflow 订单处理工作流
func OrderWorkflow(ctx workflow.Context, orderID string) error {
	// 获取带 trace 的 workflow logger
	// 自动包含 workflow_id, run_id, trace_id 等信息
	logger := logger.NewWorkflowLogger(ctx)

	logger.Info("开始处理订单", "order_id", orderID)

	// 执行 Activity
	var result string
	activityCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	})

	logger.Debug("准备执行支付 Activity")

	err := workflow.ExecuteActivity(activityCtx, ProcessPaymentActivity, orderID).Get(ctx, &result)
	if err != nil {
		logger.Error("支付处理失败", "order_id", orderID, "error", err)
		return err
	}

	logger.Info("订单处理完成", "order_id", orderID, "result", result)
	return nil
}

// ========================= Activity 中使用 =========================

// ProcessPaymentActivity 处理支付
func ProcessPaymentActivity(ctx context.Context, orderID string) (string, error) {
	// 获取带 trace 的 activity logger
	// 自动包含 activity_type, workflow_id, trace_id 等信息
	logger := logger.NewActivityLogger(ctx)

	logger.Info("开始处理支付", "order_id", orderID)

	// 模拟支付处理
	logger.Debug("调用支付网关", "gateway", "stripe")

	// 使用 WithFields 添加更多上下文
	paymentLog := logger.WithFields(map[string]interface{}{
		"payment_method": "credit_card",
		"currency":       "USD",
	})

	paymentLog.Info("支付验证中", "amount", 99.99)

	// 模拟处理时间
	time.Sleep(100 * time.Millisecond)

	paymentLog.Info("支付成功", "transaction_id", "TXN-123456")

	return "payment_successful", nil
}

// ========================= 高级配置 =========================

func AdvancedConfigExample() {
	// 完全自定义配置
	config := logger.LoggerConfig{
		Level:       "info",
		WithIcon:    true,      // 启用图标
		WithColor:   true,      // 启用颜色
		Format:      "console", // console 或 json
		ShowCaller:  true,      // 显示调用位置
		ProjectRoot: "/app",    // 设置项目根目录（用于相对路径）
		Fields: map[string]interface{}{
			"app_name": "my-service",
			"version":  "2.0.0",
		},
	}

	customLog := logger.NewLoggerWithConfig(config)
	customLog.Info("使用自定义配置的日志")
}

// ========================= 迁移指南 =========================

func MigrationExample(ctx workflow.Context, actCtx context.Context) {
	// 旧代码（仍然兼容）：
	// import "github.com/hanyougame/glib/flowcore/trace"
	// wfLog := trace.GetWorkflowLogger(ctx)
	// actLog := trace.GetActivityLogger(actCtx)

	// 新代码（推荐）：
	// import "github.com/hanyougame/glib/flowcore/logger"
	wfLog := logger.NewWorkflowLogger(ctx)
	actLog := logger.NewActivityLogger(actCtx)

	wfLog.Info("Workflow 日志")
	actLog.Info("Activity 日志")

	// 新功能：使用 WithFields 添加额外字段
	customLog := wfLog.WithFields(map[string]interface{}{
		"custom_field": "custom_value",
	})
	customLog.Info("带自定义字段的日志")
}

// ========================= 输出示例 =========================

/*
开发环境输出（带颜色和图标）：
2024-01-15 10:30:45.123 | ℹ️ [INFO] | 📂 example/usage.go:15:BasicUsage | temporal | 应用启动 | version=1.0.0 env=production

生产环境输出（JSON 格式）：
{"time":"2024-01-15 10:30:45.123","level":"INFO","caller":"@ example/usage.go:20:BasicUsage","logger":"temporal","msg":"订单创建","order_id":"ORD-123456","amount":99.99}

Workflow 日志（自动包含 trace）：
2024-01-15 10:30:45.123 | ℹ️ [INFO] | 📂 workflow/order.go:25:OrderWorkflow | temporal | 开始处理订单 | order_id=ORD-789 workflow_id=order-workflow-123 run_id=abc-def-ghi trace_id=xyz-123

Activity 日志（自动包含 trace）：
2024-01-15 10:30:45.456 | ℹ️ [INFO] | 📂 activity/payment.go:15:ProcessPayment | temporal | 支付成功 | transaction_id=TXN-123456 activity_type=ProcessPaymentActivity workflow_id=order-workflow-123 trace_id=xyz-123
*/
