package tracer

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
)

func ServerTracePropagator() stats.Handler {
	return new(serverTraceProvider)
}

type serverTraceProvider struct{}

func (s *serverTraceProvider) TagRPC(ctx context.Context, _ *stats.RPCTagInfo) context.Context {
	md, hasMD := metadata.FromIncomingContext(ctx)
	if !hasMD {
		return ctx
	}
	var err error
	config := trace.SpanContextConfig{
		TraceFlags: trace.FlagsSampled,
		Remote:     true,
	}
	traceHeaders, hasTraceHeaders := md["x-trace-id"]
	if !hasTraceHeaders || len(traceHeaders) == 0 {
		return ctx
	}
	if config.TraceID, err = trace.TraceIDFromHex(traceHeaders[0]); err != nil {
		return ctx
	}

	spanHeaders, hasSpanHeaders := md["x-span-id"]
	if !hasSpanHeaders || len(spanHeaders) == 0 {
		return ctx
	}
	if config.SpanID, err = trace.SpanIDFromHex(spanHeaders[0]); err != nil {
		return ctx
	}
	spanContext := trace.NewSpanContext(config)
	if !spanContext.IsValid() {
		return ctx
	}
	return trace.ContextWithSpanContext(ctx, spanContext)
}

func (s *serverTraceProvider) HandleRPC(_ context.Context, _ stats.RPCStats) {
}

func (s *serverTraceProvider) TagConn(ctx context.Context, _ *stats.ConnTagInfo) context.Context {
	return ctx
}

func (s *serverTraceProvider) HandleConn(_ context.Context, _ stats.ConnStats) {
}
