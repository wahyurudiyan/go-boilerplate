package telemetry

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

func newMeterProvider(ctx context.Context, interval int, resource *resource.Resource) (*sdkmetric.MeterProvider, error) {
	var (
		err      error
		exporter *otlpmetricgrpc.Exporter
	)

	if exporter, err = otlpmetricgrpc.New(ctx); err != nil {
		return nil, err
	}

	// Set default interval
	if interval <= 1 {
		interval = 1
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(
				exporter, sdkmetric.WithInterval(time.Duration(interval)*time.Second),
			),
		),
	)
	otel.SetMeterProvider(mp)

	return mp, nil
}
