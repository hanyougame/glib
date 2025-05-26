package redisx

import (
	"context"
	"github.com/hanyougame/glib/tracing"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logc"
	"go.opentelemetry.io/otel/attribute"
	oteltrace "go.opentelemetry.io/otel/trace"
	"net"
)

type DebugHook struct {
}

func (DebugHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}
func (DebugHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		logc.Debugf(ctx, "redis cmd: %s", cmd.String())
		return next(ctx, cmd)
	}
}

func (DebugHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		logc.Debugf(ctx, "redis cmd: %s", func() []string {
			var cmdsStr []string
			for _, cmd := range cmds {
				cmdsStr = append(cmdsStr, cmd.String())
			}
			return cmdsStr
		}())
		return next(ctx, cmds)
	}
}

// TraceHook redis追踪钩子
type TraceHook struct{}

func (TraceHook) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}
func (TraceHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		tracing.Inject(ctx, "redis", func(span oteltrace.Span) oteltrace.Span {
			span.SetAttributes(
				attribute.KeyValue{
					Key:   "redis.command",
					Value: attribute.StringValue(cmd.String()),
				},
			)
			return span
		})

		return next(ctx, cmd)
	}
}

func (TraceHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		tracing.Inject(ctx, "redis", func(span oteltrace.Span) oteltrace.Span {
			span.SetAttributes(
				attribute.KeyValue{
					Key: "redis.command",
					Value: attribute.StringSliceValue((func() []string {
						var cmdsStr []string
						for _, cmd := range cmds {
							cmdsStr = append(cmdsStr, cmd.String())
						}
						return cmdsStr
					})()),
				},
			)
			return span
		})

		return next(ctx, cmds)
	}
}
