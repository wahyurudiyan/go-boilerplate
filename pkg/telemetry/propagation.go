package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/propagation"
)

func newPropagator(ctx context.Context) propagation.TextMapPropagator {
	propagator := propagation.NewCompositeTextMapPropagator(propagation.Baggage{}, propagation.TraceContext{})
	return propagator
}
