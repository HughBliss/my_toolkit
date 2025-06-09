package tracer

import (
	"context"
	"fmt"
	zfg "github.com/chaindead/zerocfg"
	"github.com/hughbliss/my_toolkit/tracer/exporter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	otelsdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.32.0"
	"time"
)

var (
	Provider *otelsdk.TracerProvider

	enableJaegerExporter = zfg.Bool("enable", true, "JAEGER_ENABLE", zfg.Group(exporter.JaegerGroup))
)

func Init(ctx context.Context, appName string, appVer string, env string) (func(), error) {

	var options []otelsdk.TracerProviderOption

	options = append(options, otelsdk.WithResource(resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(appName),
		semconv.ServiceVersion(appVer),
		attribute.String("environment", env),
	)))
	if *enableJaegerExporter {
		jaegerExporter, err := exporter.Jaeger(ctx)
		if err != nil {
			return nil, err
		}
		options = append(options, otelsdk.WithBatcher(jaegerExporter,
			otelsdk.WithMaxQueueSize(otelsdk.DefaultMaxQueueSize*4),
			otelsdk.WithMaxExportBatchSize(otelsdk.DefaultMaxExportBatchSize*4),
			otelsdk.WithExportTimeout(otelsdk.DefaultExportTimeout*4*time.Millisecond),
			otelsdk.WithBatchTimeout(otelsdk.DefaultScheduleDelay*4*time.Millisecond),
		))
	}

	Provider = otelsdk.NewTracerProvider(options...)

	otel.SetTracerProvider(Provider)

	return func() {
		if err := Provider.Shutdown(ctx); err != nil {
			fmt.Printf("failed to shutdown tracer: %s\n", err.Error())
		}
	}, nil
}
