package telemetry

import (
	"context"
	"errors"
	"time"

	otelruntime "go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type ShutdownCallback func(context.Context) error

type TelemetrySetup struct {
	ServiceName        string
	ServiceVersion     string
	EnableRuntimeMeter bool
	Interval           int // in second
}

func SetupOpentelemetry(ctx context.Context, setup TelemetrySetup) (ShutdownCallback, error) {
	var (
		err               error
		shutdownCallbacks []ShutdownCallback
	)

	shutdown := func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownCallbacks {
			err = errors.Join(err, fn(ctx))
		}
		shutdownCallbacks = nil

		return err
	}

	handleErr := func(e error) {
		err = errors.Join(e, shutdown(ctx))
	}

	// Init resource
	serviceName := "go-boilerplate-service"
	if setup.ServiceName != "" {
		serviceName = setup.ServiceName
	}

	serviceVersion := "v0.0.0"
	if setup.ServiceVersion != "" {
		serviceName = setup.ServiceVersion
	}

	resource, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		),
	)
	if err != nil {
		handleErr(err)
		return nil, err
	}

	// Init propagator
	propagator := newPropagator(ctx)
	otel.SetTextMapPropagator(propagator)

	// Init tracer
	traceProvider, err := newTracerProvider(ctx, resource)
	if err != nil {
		handleErr(err)
		return nil, err
	}
	otel.SetTracerProvider(traceProvider)

	// Init meter
	meterProvider, err := newMeterProvider(ctx, setup.Interval, resource)
	if err != nil {
		handleErr(err)
		return nil, err
	}
	otel.SetMeterProvider(meterProvider)

	// Setup runtime meter to monitor golang runtime info such as goroutine count, garbage collector, etc.\
	// NOTES: This metrics is under development, please use your wisdom.
	if setup.EnableRuntimeMeter {
		if err := otelruntime.Start(
			otelruntime.WithMinimumReadMemStatsInterval(time.Duration(setup.Interval) * time.Second),
		); err != nil {
			handleErr(err)
			return nil, err
		}

	}

	// Everything should be okay, now append the shutdown for all provider
	shutdownCallbacks = append(shutdownCallbacks, traceProvider.Shutdown, meterProvider.Shutdown)

	return shutdown, nil
}
