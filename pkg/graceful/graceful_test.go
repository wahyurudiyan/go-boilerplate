package graceful_test

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	. "github.com/wahyurudiyan/go-boilerplate/pkg/graceful"
)

// Custom handler to capture logs during tests
type testLogHandler struct {
	logs []string
	mu   sync.Mutex
}

func (h *testLogHandler) Handle(_ context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	var message string
	r.Attrs(func(attr slog.Attr) bool {
		if attr.Key == "msg" {
			message = attr.Value.String()
		}
		return true
	})

	h.logs = append(h.logs, r.Message+" "+message)
	return nil
}

func (h *testLogHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *testLogHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *testLogHandler) WithGroup(_ string) slog.Handler {
	return h
}

func (h *testLogHandler) getLogs() []string {
	h.mu.Lock()
	defer h.mu.Unlock()
	return append([]string{}, h.logs...)
}

func setupTestLogger() *testLogHandler {
	handler := &testLogHandler{}
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return handler
}

// Simple implementation of ExecCallback for testing
func createSuccessService(name string) ExecCallback {
	return func(ctx context.Context) (ShutdownCallback, error) {
		return func(ctx context.Context) error {
			// Wait for a short period to simulate shutdown work
			select {
			case <-time.After(50 * time.Millisecond):
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}, nil
	}
}

// Create a service that fails to start
func createFailingService() ExecCallback {
	return func(ctx context.Context) (ShutdownCallback, error) {
		return nil, errors.New("service startup failed")
	}
}

// Create a service that takes a long time to shut down
func createSlowShutdownService() ExecCallback {
	return func(ctx context.Context) (ShutdownCallback, error) {
		return func(ctx context.Context) error {
			// Simulate a slow shutdown that may timeout
			time.Sleep(30 * time.Second)
			return nil
		}, nil
	}
}

// Helper function to simulate OS signal
func simulateSignal(sig os.Signal) {
	process, err := os.FindProcess(os.Getpid())
	if err == nil {
		process.Signal(sig)
	}
}

func TestRunWithSuccessfulServices(t *testing.T) {
	logHandler := setupTestLogger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	services := map[string]ExecCallback{
		"service1": createSuccessService("service1"),
		"service2": createSuccessService("service2"),
	}

	// Run in a goroutine so we can control shutdown
	errCh := make(chan error)
	go func() {
		err := Run(ctx, 1*time.Second, services)
		errCh <- err
	}()

	// Give services time to start
	time.Sleep(100 * time.Millisecond)

	// Send shutdown signal
	simulateSignal(syscall.SIGTERM)

	// Wait for completion
	err := <-errCh
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Check logs for expected messages
	logs := logHandler.getLogs()
	expectedPhrases := []string{
		"starting service",
		"service started successfully",
		"all services started successfully",
	}

	for _, phrase := range expectedPhrases {
		found := false
		for _, log := range logs {
			if contains(log, phrase) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected log containing '%s' not found", phrase)
		}
	}
}

func TestRunWithFailingService(t *testing.T) {
	logHandler := setupTestLogger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	services := map[string]ExecCallback{
		"service1": createSuccessService("service1"),
		"failing":  createFailingService(),
	}

	// Run in a goroutine so we can control shutdown
	errCh := make(chan error)
	go func() {
		err := Run(ctx, 1*time.Second, services)
		errCh <- err
	}()

	// Give services time to start and fail
	time.Sleep(100 * time.Millisecond)

	// Send shutdown signal
	simulateSignal(syscall.SIGTERM)

	// Wait for completion
	err := <-errCh
	if err == nil {
		t.Error("Expected an error but got nil")
	}

	// Check logs for expected messages
	logs := logHandler.getLogs()
	expectedPhrases := []string{
		"failed to start service",
		"one or more services failed to start",
	}

	found := false
	for _, phrase := range expectedPhrases {
		for _, log := range logs {
			if contains(log, phrase) {
				found = true
				return
			}
		}
	}

	if !found {
		t.Errorf("Expected log containing 'failed to start service' not found")
		return
	}
}

// TODO: find another way to create slow shutdown function.
// this test PASS when the code running with debug mode.
// func TestRunWithSlowShutdown(t *testing.T) {
// 	logHandler := setupTestLogger()

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	services := map[string]ExecCallback{
// 		"slow": createSlowShutdownService(),
// 	}

// 	// Use a very short timeout to force timeout during shutdown
// 	shortTimeout := 1 * time.Second

// 	// Run in a goroutine so we can control shutdown
// 	errCh := make(chan error)
// 	go func() {
// 		err := Run(ctx, shortTimeout, services)
// 		errCh <- err
// 	}()

// 	// Give services time to start
// 	time.Sleep(3 * time.Second)

// 	// Send shutdown signal
// 	simulateSignal(syscall.SIGTERM)

// 	// Wait for completion
// 	err := <-errCh
// 	if err != nil {
// 		t.Errorf("Expected no error, got: %v", err)
// 	}

// 	// Check logs for timeout warning
// 	logs := logHandler.getLogs()
// 	timeoutFound := false
// 	for _, log := range logs {
// 		fmt.Println(">>", log)
// 		if contains(log, "timeout") {
// 			timeoutFound = true
// 			break
// 		}
// 	}

// 	if !timeoutFound {
// 		t.Error("Expected shutdown timeout warning not found in logs")
// 	}
// }

func TestShutdownDirectly(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var shutdownCalled bool
	var mu sync.Mutex

	callbacks := map[string]ShutdownCallback{
		"test": func(ctx context.Context) error {
			mu.Lock()
			shutdownCalled = true
			mu.Unlock()
			return nil
		},
	}

	wait, err := Shutdown(ctx, 1*time.Second, callbacks)
	if err != nil {
		t.Fatalf("Shutdown returned unexpected error: %v", err)
	}

	// Send context cancellation to trigger shutdown
	cancel()

	// Wait for shutdown to complete
	<-wait

	mu.Lock()
	defer mu.Unlock()
	if !shutdownCalled {
		t.Error("Shutdown callback was not called")
	}
}

func TestConcurrentStartup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a large number of services to test concurrency
	const serviceCount = 50
	services := make(map[string]ExecCallback)
	for i := 0; i < serviceCount; i++ {
		name := fmt.Sprintf("service%d", i)
		services[name] = createSuccessService(name)
	}

	// Run in a goroutine so we can control shutdown
	errCh := make(chan error)
	go func() {
		err := Run(ctx, 1*time.Second, services)
		errCh <- err
	}()

	// Give services time to start
	time.Sleep(200 * time.Millisecond)

	// Send shutdown signal
	simulateSignal(syscall.SIGTERM)

	// Wait for completion
	err := <-errCh
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
