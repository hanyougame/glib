package pkg

import (
	"sync"
	"time"

	"go.temporal.io/sdk/temporal"
)

// 预定义的重试策略名称
const (
	Standard = "standard" // 标准重试策略
	Critical = "critical" // 关键业务重试策略
	Minimal  = "minimal"  // 最小化重试策略
	NoRetry  = "none"     // 不重试策略
)

var (
	mu       sync.RWMutex
	policies = make(map[string]*temporal.RetryPolicy) // 存储重试策略的映射表
)

// Register 添加或更新一个重试策略
func Register(name string, policy *temporal.RetryPolicy) {
	mu.Lock()
	defer mu.Unlock()
	policies[name] = policy
}

// Get 根据名称获取一个重试策略，如果未找到则返回 standard 策略
func Get(name string) *temporal.RetryPolicy {
	mu.RLock()
	defer mu.RUnlock()

	if policy, ok := policies[name]; ok {
		return policy
	}

	// 如果未找到，则返回标准策略
	if standard, ok := policies[Standard]; ok {
		return standard
	}

	return nil
}

// InitializeDefaultPolicies 初始化默认的重试策略
func InitializeDefaultPolicies() {
	// 标准重试策略
	Register(Standard, &temporal.RetryPolicy{
		InitialInterval:    time.Second, // 初始重试间隔时间
		BackoffCoefficient: 2.0,         // 退避系数
		MaximumInterval:    time.Minute, // 最大重试间隔时间
		MaximumAttempts:    5,           // 最大重试次数
	})

	// 关键业务重试策略
	Register(Critical, &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    5 * time.Minute, // 最大重试间隔时间为 5 分钟
		MaximumAttempts:    20,              // 最大重试次数为 20 次
	})

	// 最小化重试策略
	Register(Minimal, &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    time.Minute,
		MaximumAttempts:    2, // 最大重试次数为 2 次
	})

	// 不重试策略
	Register(NoRetry, &temporal.RetryPolicy{
		MaximumAttempts: 1, // 最大重试次数为 1，即不重试
	})
}

// WithNonRetryableErrors 向策略中添加不可重试的错误类型
func WithNonRetryableErrors(policy *temporal.RetryPolicy, errorTypes ...string) *temporal.RetryPolicy {
	if policy == nil {
		return nil
	}

	newPolicy := *policy
	if len(policy.NonRetryableErrorTypes) > 0 {
		// 组合已有的和新的不可重试错误类型
		combined := make([]string, len(policy.NonRetryableErrorTypes)+len(errorTypes))
		copy(combined, policy.NonRetryableErrorTypes)
		copy(combined[len(policy.NonRetryableErrorTypes):], errorTypes)
		newPolicy.NonRetryableErrorTypes = combined
	} else {
		// 直接使用传入的错误类型作为不可重试类型
		newPolicy.NonRetryableErrorTypes = errorTypes
	}

	return &newPolicy
}
