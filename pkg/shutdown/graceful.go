package shutdown

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

type FnOps func(ctx context.Context) error

func Graceful(ctx context.Context, timeout time.Duration, ops map[string]FnOps) {
	/**
	Define graceful shutdown
	*/
	shutdownSignals := []os.Signal{
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	}
	ctx, stop := signal.NotifyContext(ctx, shutdownSignals...)
	defer stop()

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(len(ops))
	for name, operation := range ops {
		go func(ctx context.Context) {
			defer wg.Done()
			slog.InfoContext(ctx, fmt.Sprintf("%s is shutting down", name))
			if err := operation(ctx); err != nil {
				slog.ErrorContext(ctx, fmt.Sprintf("%s is failed to shutting down", name), "error", err)
				return
			}
		}(ctx)
	}
}
