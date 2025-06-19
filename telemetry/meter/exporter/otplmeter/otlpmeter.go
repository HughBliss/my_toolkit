package otplmeterexporter

import (
	"context"
	"fmt"
	zfg "github.com/chaindead/zerocfg"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding/gzip"
)

var (
	otlpmeterGroup = zfg.NewGroup("otlpmeter")
	otlpmeterHost  = zfg.Str("host", "0.0.0.0", "OTLPMETER_HOST", zfg.Group(otlpmeterGroup))
	otlpmeterPort  = zfg.Uint32("port", 4317, "OTLPMETER_HOST", zfg.Group(otlpmeterGroup))
)

func OTLPMeter(ctx context.Context) (metric.Reader, error) {
	endpoint := fmt.Sprintf("%s:%d", *otlpmeterHost, *otlpmeterPort)
	fmt.Printf("exporting metrics to otlp meter at %s\n", endpoint)
	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(endpoint),
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithCompressor(gzip.Name),
		otlpmetricgrpc.WithDialOption(grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*10))),
	)
	if err != nil {
		return nil, err
	}

	return metric.NewPeriodicReader(exporter), nil
}
