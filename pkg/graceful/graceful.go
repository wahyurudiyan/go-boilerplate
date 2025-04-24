package graceful

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type ShutdownCallback func(ctx context.Context) error

type ExecCallback func(ctx context.Context) (ShutdownCallback, error)

func Runner(ctx context.Context, timeout time.Duration, ops map[string]ExecCallback) (<-chan struct{}, error) {
	var (
		shutdownCb = make(map[string]ShutdownCallback)
	)

	// Execute operations and collect shutdown callbacks
	for opName, opFn := range ops {
		sdCallback, err := opFn(ctx)
		if err != nil {
			// Clean up any already started operations
			cleanupCtx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()

			for name, callback := range shutdownCb {
				if err := callback(cleanupCtx); err != nil {
					slog.ErrorContext(ctx, fmt.Sprintf("%s failed during cleanup", name), "error", err)
				}
			}

			return nil, fmt.Errorf("operation %s failed to start: %w", opName, err)
		}

		shutdownCb[opName] = sdCallback
	}

	// Register shutdown handlers
	wait, err := Shutdown(ctx, timeout, shutdownCb)
	if err != nil {
		return nil, err
	}

	return wait, nil
}

func Shutdown(ctx context.Context, timeout time.Duration, ops map[string]ShutdownCallback) (<-chan struct{}, error) {
	wait := make(chan struct{})

	go func() {
		// Define graceful shutdown
		shutdownSignals := []os.Signal{
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
		}

		sigCtx, stop := signal.NotifyContext(ctx, shutdownSignals...)
		defer stop()

		// Wait for signal
		<-sigCtx.Done()
		slog.InfoContext(ctx, "Shutdown signal received, initiating graceful shutdown")

		// Create timeout context for shutdown operations
		shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		var wg sync.WaitGroup

		wg.Add(len(ops))
		for name, operation := range ops {
			// Capture loop variables to avoid data race
			opName := name
			opFunc := operation

			go func() {
				defer wg.Done()
				slog.InfoContext(shutdownCtx, fmt.Sprintf("%s is shutting down", opName))
				if err := opFunc(shutdownCtx); err != nil {
					slog.ErrorContext(shutdownCtx, fmt.Sprintf("%s failed to shut down", opName), "error", err)
					return
				}
				slog.InfoContext(shutdownCtx, fmt.Sprintf("%s shut down successfully", opName))
			}()
		}

		// Wait for all operations to complete or timeout
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			slog.InfoContext(ctx, "All services shut down successfully")
		case <-shutdownCtx.Done():
			slog.WarnContext(ctx, "Shutdown timed out before all services could shut down cleanly")
		}

		close(wait)
	}()

	return wait, nil
}
