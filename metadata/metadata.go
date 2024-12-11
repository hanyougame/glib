package metadata

import (
	"context"
	"github.com/spf13/cast"
)

const (
	CtxJWTUserId   = "uid"
	CtxJWTUsername = "username"
)

// WithMetadata 上下文数据
func WithMetadata(ctx context.Context, key, val any) context.Context {
	return context.WithValue(ctx, key, val)
}

// GetMetadataFromCtx 获取上下文数据
func GetMetadataFromCtx(ctx context.Context, key any) any {
	return ctx.Value(key)
}

// GetMetadata 上下文取值
func GetMetadata[T any](ctx context.Context, key any) (T, bool) {
	val := ctx.Value(key)
	if val == nil {
		var zero T
		return zero, false
	}

	// 断言类型
	if result, ok := val.(T); ok {
		return result, true
	}

	var zero T
	return zero, false
}

// GetUidFromCtx 从上下文中获取uid
func GetUidFromCtx(ctx context.Context) int {
	return cast.ToInt(ctx.Value(CtxJWTUserId))
}

// GetUsernameFromCtx 从上下文中获取username
func GetUsernameFromCtx(ctx context.Context) string {
	return cast.ToString(ctx.Value(CtxJWTUsername))
}
