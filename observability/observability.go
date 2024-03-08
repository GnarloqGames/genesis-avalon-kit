package observability

import (
	"context"
	"errors"

	"github.com/GnarloqGames/genesis-avalon-kit/observability/metric"
	"github.com/GnarloqGames/genesis-avalon-kit/observability/trace"
	"go.opentelemetry.io/otel"
)

func Setup(ctx context.Context, serviceName, serviceVersion string) (shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up meter provider.
	meterProvider, err := metric.New(serviceName, serviceVersion)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	traceProvider, err := trace.New(serviceName, serviceVersion)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, traceProvider.Shutdown)
	otel.SetTracerProvider(traceProvider)

	return
}
