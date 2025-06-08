package exporter

import (
	"context"
	"fmt"
	zfg "github.com/chaindead/zerocfg"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
)

var (
	jaegerGroup = zfg.NewGroup("jaeger")
	jaegerHost  = zfg.Str("host", "0.0.0.0", "JAEGER_HOST", zfg.Group(jaegerGroup))
	jaegerPort  = zfg.Uint32("port", 4317, "JAEGER_HOST", zfg.Group(jaegerGroup))
)

func Jaeger(ctx context.Context) (*otlptrace.Exporter, error) {
	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(fmt.Sprintf("%s:%d", *jaegerHost, *jaegerPort)),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	return exporter, nil
}
