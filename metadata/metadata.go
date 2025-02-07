package metadata

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"net"
)

const (
	// CtxJWTUserId   用户id
	CtxJWTUserId = "uid"
	// CtxJWTUsername 用户名
	CtxJWTUsername = "username"
	// CtxIp          ip
	CtxIp = "ip"
	// CtxDomain      域名
	CtxDomain = "domain"
	// CtxRegion       区域
	CtxRegion = "region"
	// CtxDeviceID     设备id
	CtxDeviceID = "device_id"
	// CtxDeviceType  设备类型
	CtxDeviceType = "device_type"
	// CtxBrowserFingerprint 浏览器指纹
	CtxBrowserFingerprint = "browser_fingerprint"
	// CtxCurrencyCode 币种code
	CtxCurrencyCode = "currency_code"
)

const (
	// RegionKey 地区
	RegionKey = "X-Region"
	// DeviceIDKey 设备id
	DeviceIDKey = "X-Device-ID"
	// DeviceTypeKey 设备类型
	DeviceTypeKey = "X-Device-Type"
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
	if val, ok := ctx.Value(key).(T); ok {
		return val, true
	}

	var zero T
	return zero, false
}

// GetUidFromCtx 从上下文中获取uid
func GetUidFromCtx(ctx context.Context) int64 {
	// 获取值
	val := ctx.Value(CtxJWTUserId)
	if val == nil {
		// 如果为空，可以根据需要返回默认值
		return 0
	}

	// 如果值是 json.Number 类型
	if num, ok := val.(json.Number); ok {
		// 转换为 int64
		uid, err := num.Int64()
		if err != nil {
			// 如果转换失败，可以选择处理错误（如返回 0 或日志记录）
			fmt.Println("Error converting json.Number to int64:", err)
			return 0
		}
		return uid
	}

	// 如果值不是 json.Number 类型，直接转换
	return cast.ToInt64(val)
}

// GetUsernameFromCtx 从上下文中获取username
func GetUsernameFromCtx(ctx context.Context) string {
	return cast.ToString(ctx.Value(CtxJWTUsername))
}

// GetCurrencyCodeFromCtx 从上下文中获取currency_code
func GetCurrencyCodeFromCtx(ctx context.Context) string {
	return cast.ToString(ctx.Value(CtxCurrencyCode))
}

// GetIpFromCtx 从上下文中获取ip
func GetIpFromCtx(ctx context.Context) string {
	if val := ctx.Value(CtxIp); val != nil {
		switch v := val.(type) {
		case string:
			return v
		case net.IP:
			return v.String()
		default:
			if s, ok := val.(fmt.Stringer); ok {
				return s.String()
			}
		}
	}
	return ""
}

// GetDomainFromCtx 从上下文中获取域名
func GetDomainFromCtx(ctx context.Context) string {
	if domain, ok := GetMetadata[string](ctx, CtxDomain); ok {
		return domain
	}
	return ""
}

// GetDeviceIDFromCtx 从上下文中获取设备id
func GetDeviceIDFromCtx(ctx context.Context) string {
	if deviceID, ok := GetMetadata[string](ctx, CtxDeviceID); ok {
		return deviceID
	}
	return ""
}

// GetDeviceTypeFromCtx 从上下文中获取设备类型
func GetDeviceTypeFromCtx(ctx context.Context) string {
	if deviceType, ok := GetMetadata[string](ctx, CtxDeviceType); ok {
		return deviceType
	}
	return ""
}

// GetBrowserFingerprintFromCtx 从上下文中获取浏览器指纹
func GetBrowserFingerprintFromCtx(ctx context.Context) string {
	if browserFingerprint, ok := GetMetadata[string](ctx, CtxBrowserFingerprint); ok {
		return browserFingerprint
	}
	return ""
}

// GetRegionFromCtx 从上下文中获取区域
func GetRegionFromCtx(ctx context.Context) string {
	if region, ok := GetMetadata[string](ctx, CtxRegion); ok {
		return region
	}
	return ""
}
