package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func newTracerProvider(ctx context.Context, resource *resource.Resource) (*sdktrace.TracerProvider, error) {
	var (
		err      error
		exporter *otlptrace.Exporter
	)

	if exporter, err = otlptracegrpc.New(ctx); err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
	)

	return tp, nil
}
