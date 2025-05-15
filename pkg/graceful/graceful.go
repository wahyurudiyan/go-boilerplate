package graceful

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type ShutdownCallback func(ctx context.Context) error

type ExecCallback func(ctx context.Context) (ShutdownCallback, error)

func Run(ctx context.Context, timeout time.Duration, ops map[string]ExecCallback) error {
	// Create a base context for running operations
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Start all operations and collect shutdown callbacks
	var shutdownCallbacksMu sync.Mutex
	shutdownCallbacks := make(map[string]ShutdownCallback)
	var startupWg sync.WaitGroup
	var startupErr error
	var startupErrMu sync.Mutex

	// Start all services concurrently
	for name, operation := range ops {
		startupWg.Add(1)
		go func(name string, operation ExecCallback) {
			defer startupWg.Done()

			slog.InfoContext(ctx, "[Graceful] 🚀 starting service", "name", name)

			callback, err := operation(runCtx)
			if err != nil {
				slog.ErrorContext(ctx, "[Graceful] ⛔ failed to start service", "name", name, "error", err)
				startupErrMu.Lock()
				if startupErr == nil {
					startupErr = err
				}
				startupErrMu.Unlock()
				cancel() // Cancel context to signal other operations to terminate
				return
			}

			if callback != nil {
				shutdownCallbacksMu.Lock()
				shutdownCallbacks[name] = callback
				shutdownCallbacksMu.Unlock()
				slog.InfoContext(ctx, "[Graceful] ✅ service started successfully", "name", name)
			}
		}(name, operation)
	}

	go func() {
		// Wait for all services to start or fail
		startupWg.Wait()
	}()

	// Check if any service failed to start
	startupErrMu.Lock()
	hasStartupError := startupErr != nil
	startupErrMu.Unlock()

	if hasStartupError {
		slog.ErrorContext(ctx, "[Graceful] ⛔ one or more services failed to start")
		// Still proceed with shutdown for any services that did start
	} else {
		slog.InfoContext(ctx, "[Graceful] 🌟 all services started successfully")
	}

	// Initiate shutdown process
	slog.InfoContext(ctx, "[Graceful] 🌟 service shutting down")
	shutdownCallbacksMu.Lock()
	cloneCallback := make(map[string]ShutdownCallback, len(shutdownCallbacks))
	for k, v := range shutdownCallbacks {
		cloneCallback[k] = v
	}
	shutdownCallbacksMu.Unlock()

	wait, err := Shutdown(ctx, timeout, cloneCallback)
	if err != nil {
		return err
	}

	// Wait for shutdown to complete
	<-wait

	startupErrMu.Lock()
	defer startupErrMu.Unlock()
	return startupErr
}

func Shutdown(ctx context.Context, timeout time.Duration, ops map[string]ShutdownCallback) (<-chan struct{}, error) {
	wait := make(chan struct{})

	go func() {
		defer close(wait)

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
		slog.InfoContext(ctx, "[Graceful] 📡 shutdown signal received, initiating graceful shutdown")

		// Create timeout context for shutdown operations
		shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		var wg sync.WaitGroup
		for name, operation := range ops {
			wg.Add(1)
			go func(name string, operation ShutdownCallback) {
				defer wg.Done()
				slog.Info("[Graceful] 🔻 shutting down", "name", name)
				if err := operation(shutdownCtx); err != nil {
					slog.ErrorContext(shutdownCtx, "[Graceful] ⛔ failed to shut down", "name", name, "error", err)
					return
				}
				slog.InfoContext(shutdownCtx, "[Graceful] 👍 shutdown successfully", "name", name)
			}(name, operation)
		}

		// Wait for all operations to complete or timeout
		done := make(chan struct{})
		go func() {
			wg.Wait()
			close(done)
		}()

		select {
		case <-done:
			slog.InfoContext(ctx, "[Graceful] 🛬 all services shut down successfully")
		case <-shutdownCtx.Done():
			slog.WarnContext(ctx, "[Graceful] ⏱️ shutdown timeout before all services could shut down cleanly")
		}
	}()

	return wait, nil
}
