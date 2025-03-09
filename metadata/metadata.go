package metadata

import (
	"context"
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
	// CtxChannelID 渠道id
	CtxChannelID = "channel_id"
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
	return cast.ToInt64(ctx.Value(CtxJWTUserId))
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

// GetChannelIDFromCtx 从上下文获取渠道id
func GetChannelIDFromCtx(ctx context.Context) int64 {
	if channelID, ok := GetMetadata[int64](ctx, CtxChannelID); ok {
		return channelID
	}
	return 0
}
