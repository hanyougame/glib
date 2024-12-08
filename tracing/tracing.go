package tracing

import (
	"context"
	"github.com/zeromicro/go-zero/core/trace"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type Handle func(span oteltrace.Span) oteltrace.Span

// Inject 链路追踪
func Inject(ctx context.Context, spanName string, fn Handle, carriers ...propagation.TextMapCarrier) {
	tracer := otel.GetTracerProvider().Tracer(trace.TraceName)
	spanCtx, span := tracer.Start(ctx, spanName, oteltrace.WithSpanKind(oteltrace.SpanKindProducer))
	defer span.End()

	var carrier propagation.TextMapCarrier = &propagation.HeaderCarrier{}
	if len(carriers) > 0 {
		carrier = carriers[0]
	}
	span = fn(span)
	otel.GetTextMapPropagator().Inject(spanCtx, carrier)
}
