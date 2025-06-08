package tracer

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
	otelsdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.32.0"
)

var (
	Provider *otelsdk.TracerProvider
)

func Init(ctx context.Context, appName string, appVer string, exporters ...otelsdk.SpanExporter) func() {

	var options []otelsdk.TracerProviderOption

	options = append(options, otelsdk.WithResource(resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(appName),
		semconv.ServiceVersion(appVer),
		attribute.String("environment", "debug"),
	)))

	for _, exporter := range exporters {
		options = append(options, otelsdk.WithBatcher(exporter))
	}

	Provider = otelsdk.NewTracerProvider(options...)

	otel.SetTracerProvider(Provider)

	return func() {
		if err := Provider.Shutdown(ctx); err != nil {
			fmt.Printf("failed to shutdown tracer: %s\n", err.Error())
		}
	}
}
