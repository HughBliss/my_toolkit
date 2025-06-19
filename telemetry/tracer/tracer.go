package tracer

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	otelsdk "go.opentelemetry.io/otel/sdk/trace"
	"time"
)

func Init(ctx context.Context, r *resource.Resource, exporters ...otelsdk.SpanExporter) func() {

	var options []otelsdk.TracerProviderOption

	options = append(options, otelsdk.WithResource(r))

	for _, exp := range exporters {
		options = append(options, otelsdk.WithBatcher(exp,
			otelsdk.WithMaxQueueSize(otelsdk.DefaultMaxQueueSize*4),
			otelsdk.WithMaxExportBatchSize(otelsdk.DefaultMaxExportBatchSize*4),
			otelsdk.WithExportTimeout(otelsdk.DefaultExportTimeout*4*time.Millisecond),
			otelsdk.WithBatchTimeout(otelsdk.DefaultScheduleDelay*4*time.Millisecond),
		))
	}

	provider := otelsdk.NewTracerProvider(options...)

	otel.SetTracerProvider(provider)

	return func() {
		if err := provider.Shutdown(ctx); err != nil {
			fmt.Printf("failed to shutdown tracer: %s\n", err.Error())
		}
	}
}
