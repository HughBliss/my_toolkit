package otpltraceexporter

import (
	"context"
	"fmt"
	zfg "github.com/chaindead/zerocfg"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

var (
	otlptraceGroup = zfg.NewGroup("otlptrace")
	otlptraceHost  = zfg.Str("host", "0.0.0.0", "OTLPTRACE_HOST", zfg.Group(otlptraceGroup))
	otlptracePort  = zfg.Uint32("port", 4317, "OTLPTRACE_HOST", zfg.Group(otlptraceGroup))
)

func OTLPTrace(ctx context.Context) (*otlptrace.Exporter, error) {
	return otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(fmt.Sprintf("%s:%d", *otlptraceHost, *otlptracePort)),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithCompressor(gzip.Name),
		otlptracegrpc.WithDialOption(grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*10))),
	)
}
