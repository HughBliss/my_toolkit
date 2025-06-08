package tracer

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/stats"
)

type clientTraceProvider struct{}

func ClientTraceProvider() stats.Handler {
	return new(clientTraceProvider)
}

func (t *clientTraceProvider) TagRPC(ctx context.Context, _ *stats.RPCTagInfo) context.Context {
	spanCtx := trace.SpanContextFromContext(ctx)
	if !spanCtx.IsValid() {
		return ctx
	}
	if spanCtx.HasTraceID() {
		ctx = metadata.AppendToOutgoingContext(ctx, "x-trace-id", spanCtx.TraceID().String())
	}
	if spanCtx.HasSpanID() {
		ctx = metadata.AppendToOutgoingContext(ctx, "x-span-id", spanCtx.SpanID().String())
	}
	return ctx
}

func (t *clientTraceProvider) HandleRPC(_ context.Context, _ stats.RPCStats) {
}

func (t *clientTraceProvider) TagConn(ctx context.Context, _ *stats.ConnTagInfo) context.Context {
	return ctx
}

func (t *clientTraceProvider) HandleConn(_ context.Context, _ stats.ConnStats) {
}
