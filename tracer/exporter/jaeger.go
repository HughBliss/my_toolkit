package exporter

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
	jaegerGroup = zfg.NewGroup("jaeger")
	jaegerHost  = zfg.Str("host", "0.0.0.0", "JAEGER_HOST", zfg.Group(jaegerGroup))
	jaegerPort  = zfg.Uint32("port", 4317, "JAEGER_HOST", zfg.Group(jaegerGroup))
)

func Jaeger(ctx context.Context) (*otlptrace.Exporter, error) {
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(fmt.Sprintf("%s:%d", *jaegerHost, *jaegerPort)),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithCompressor(gzip.Name),
		otlptracegrpc.WithDialOption(grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*10))),
	)
	if err != nil {
		return nil, err
	}
	return exporter, nil
}
