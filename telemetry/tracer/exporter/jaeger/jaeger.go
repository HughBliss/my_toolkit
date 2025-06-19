package jaegerexporter

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
	JaegerGroup = zfg.NewGroup("jaeger")
	jaegerHost  = zfg.Str("host", "0.0.0.0", "JAEGER_HOST", zfg.Group(JaegerGroup))
	jaegerPort  = zfg.Uint32("port", 4317, "JAEGER_HOST", zfg.Group(JaegerGroup))
)

func Jaeger(ctx context.Context) (*otlptrace.Exporter, error) {
	endpoint := fmt.Sprintf("%s:%d", *jaegerHost, *jaegerPort)
	fmt.Printf("exporting traces to jaeger at %s\n", endpoint)
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithCompressor(gzip.Name),
		otlptracegrpc.WithDialOption(grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*10))),
	)
	if err != nil {
		return nil, err
	}
	return exporter, nil
}
