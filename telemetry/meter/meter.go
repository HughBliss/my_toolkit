package meter

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
	otelsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

func Init(ctx context.Context, r *resource.Resource, readers ...otelsdk.Reader) func() {

	var options []otelsdk.Option

	options = append(options, otelsdk.WithResource(r))

	for _, rdr := range readers {
		options = append(options, otelsdk.WithReader(rdr))
	}

	provider := otelsdk.NewMeterProvider(options...)

	otel.SetMeterProvider(provider)

	return func() {
		if err := provider.Shutdown(ctx); err != nil {
			fmt.Printf("failed to shutdown meter: %s\n", err.Error())
		}
	}
}
