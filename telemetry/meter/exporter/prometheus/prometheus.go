package prometheusexporter

import (
	"context"
	"fmt"
	zfg "github.com/chaindead/zerocfg"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
)

var (
	prometheusGroup = zfg.NewGroup("prometheus")
	prometheusHost  = zfg.Str("host", "0.0.0.0", "PROMETHEUS_HOST", zfg.Group(prometheusGroup))
	prometheusPort  = zfg.Uint32("port", 9090, "PROMETHEUS_HOST", zfg.Group(prometheusGroup))
)

func Prometheus(ctx context.Context) (metric.Reader, error) {
	endpoint := fmt.Sprintf("%s:%d", *prometheusHost, *prometheusPort)
	fmt.Printf("exporting metrics to prometheus at %s\n", endpoint)

	exporter, err := prometheus.New(
		prometheus.WithoutTargetInfo(),
	)
	if err != nil {
		return nil, err
	}

	return exporter, nil

}
