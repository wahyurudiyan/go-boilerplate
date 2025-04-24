package graceful

import (
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"
)

// mockOperation is a helper to create a mock operation for testing
type mockOperation struct {
	name           string
	execDelay      time.Duration // How long execution takes
	shutdownDelay  time.Duration // How long shutdown takes
	execError      error         // Error to return on execution
	shutdownError  error         // Error to return on shutdown
	shutdownCalled bool          // Whether shutdown was called
	execCalled     bool          // Whether exec was called
}

func (m *mockOperation) exec(ctx context.Context) (ShutdownCallback, error) {
	m.execCalled = true

	if m.execDelay > 0 {
		select {
		case <-time.After(m.execDelay):
		case <-ctx.Done():
			return nil, fmt.Errorf("context canceled before execution completed: %w", ctx.Err())
		}
	}

	if m.execError != nil {
		return nil, m.execError
	}

	return m.shutdown, nil
}

func (m *mockOperation) shutdown(ctx context.Context) error {
	m.shutdownCalled = true

	if m.shutdownDelay > 0 {
		select {
		case <-time.After(m.shutdownDelay):
		case <-ctx.Done():
			return fmt.Errorf("shutdown timed out: %w", ctx.Err())
		}
	}

	return m.shutdownError
}

func TestRunner_Success(t *testing.T) {
	// Setup test operations
	op1 := &mockOperation{name: "op1"}
	op2 := &mockOperation{name: "op2"}

	operations := map[string]ExecCallback{
		"op1": op1.exec,
		"op2": op2.exec,
	}

	ctx := context.Background()

	// Run with a generous timeout
	_, err := Runner(ctx, 5*time.Second, operations)
	if err != nil {
		t.Fatalf("Runner failed with error: %v", err)
	}

	// Verify that exec was called for both operations
	if !op1.execCalled {
		t.Error("op1.exec was not called")
	}
	if !op2.execCalled {
		t.Error("op2.exec was not called")
	}
}

func TestRunner_ExecError(t *testing.T) {
	// Setup test operations with one failing
	op1 := &mockOperation{name: "op1"}
	op2 := &mockOperation{name: "op2", execError: errors.New("exec error")}

	operations := map[string]ExecCallback{
		"op1": op1.exec,
		"op2": op2.exec,
	}

	ctx := context.Background()

	// Run with a generous timeout
	_, err := Runner(ctx, 5*time.Second, operations)
	if err == nil {
		t.Fatal("Runner should have failed but didn't")
	}

	// Verify that exec was called and shutdown was called for cleanup
	if !op1.execCalled {
		t.Error("op1.exec was not called")
	}
	if !op2.execCalled {
		t.Error("op2.exec was not called")
	}

	// op1 should have its shutdown called since it succeeded and needs cleanup
	if !op1.shutdownCalled {
		t.Error("op1.shutdown was not called for cleanup after op2 failed")
	}
}

func TestShutdown_AllSucceed(t *testing.T) {
	// Setup test operations
	op1 := &mockOperation{name: "op1"}
	op2 := &mockOperation{name: "op2"}

	operations := map[string]ShutdownCallback{
		"op1": op1.shutdown,
		"op2": op2.shutdown,
	}

	ctx := context.Background()

	// Create a cancelable context to simulate a signal
	cancelCtx, cancel := context.WithCancel(ctx)

	// Run shutdown
	wait, err := Shutdown(cancelCtx, 5*time.Second, operations)
	if err != nil {
		t.Fatalf("Shutdown setup failed with error: %v", err)
	}

	// Trigger shutdown signal
	cancel()

	// Wait for completion with timeout
	select {
	case <-wait:
		// Success
	case <-time.After(6 * time.Second):
		t.Fatal("Shutdown did not complete within expected time")
	}

	// Verify shutdown was called for both
	if !op1.shutdownCalled {
		t.Error("op1.shutdown was not called")
	}
	if !op2.shutdownCalled {
		t.Error("op2.shutdown was not called")
	}
}

func TestShutdown_WithError(t *testing.T) {
	// Setup test operations with one returning error
	op1 := &mockOperation{name: "op1"}
	op2 := &mockOperation{name: "op2", shutdownError: errors.New("shutdown error")}

	operations := map[string]ShutdownCallback{
		"op1": op1.shutdown,
		"op2": op2.shutdown,
	}

	ctx := context.Background()

	// Create a cancelable context to simulate a signal
	cancelCtx, cancel := context.WithCancel(ctx)

	// Run shutdown
	wait, err := Shutdown(cancelCtx, 5*time.Second, operations)
	if err != nil {
		t.Fatalf("Shutdown setup failed with error: %v", err)
	}

	// Trigger shutdown signal
	cancel()

	// Wait for completion with timeout
	select {
	case <-wait:
		// Success - should still complete even with errors
	case <-time.After(6 * time.Second):
		t.Fatal("Shutdown did not complete within expected time")
	}

	// Verify shutdown was called for both
	if !op1.shutdownCalled {
		t.Error("op1.shutdown was not called")
	}
	if !op2.shutdownCalled {
		t.Error("op2.shutdown was not called")
	}
}

func TestShutdown_Timeout(t *testing.T) {
	// Setup an operation that takes too long to shut down
	slowOp := &mockOperation{
		name:          "slow",
		shutdownDelay: 2 * time.Second,
	}

	operations := map[string]ShutdownCallback{
		"slow": slowOp.shutdown,
	}

	ctx := context.Background()

	// Create a cancelable context to simulate a signal
	cancelCtx, cancel := context.WithCancel(ctx)

	// Run shutdown with a short timeout
	wait, err := Shutdown(cancelCtx, 500*time.Millisecond, operations)
	if err != nil {
		t.Fatalf("Shutdown setup failed with error: %v", err)
	}

	// Trigger shutdown signal
	cancel()

	// Wait for completion with timeout
	select {
	case <-wait:
		// Success - should still complete even with timeout
	case <-time.After(3 * time.Second):
		t.Fatal("Shutdown did not complete despite timeout")
	}

	// Verify shutdown was called
	if !slowOp.shutdownCalled {
		t.Error("slowOp.shutdown was not called")
	}
}

func TestIntegration_SignalHandling(t *testing.T) {
	// Skip in CI environments where signal handling might be problematic
	if os.Getenv("CI") != "" {
		t.Skip("Skipping signal handling test in CI environment")
	}

	// Setup test operations
	op1 := &mockOperation{name: "op1", shutdownDelay: 100 * time.Millisecond}
	op2 := &mockOperation{name: "op2", shutdownDelay: 100 * time.Millisecond}

	operations := map[string]ShutdownCallback{
		"op1": op1.shutdown,
		"op2": op2.shutdown,
	}

	// Run shutdown with real signal handling
	wait, err := Shutdown(context.Background(), 5*time.Second, operations)
	if err != nil {
		t.Fatalf("Shutdown setup failed with error: %v", err)
	}

	// Send ourselves a signal
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		t.Fatalf("Failed to find process: %v", err)
	}

	// Use a goroutine to send the signal after a short delay
	go func() {
		time.Sleep(200 * time.Millisecond)
		err := p.Signal(syscall.SIGTERM)
		if err != nil {
			t.Errorf("Failed to send signal: %v", err)
		}
	}()

	// Wait for completion with timeout
	select {
	case <-wait:
		// Success
	case <-time.After(5 * time.Second):
		t.Fatal("Shutdown did not complete within expected time")
	}

	// Verify shutdown was called for both
	if !op1.shutdownCalled {
		t.Error("op1.shutdown was not called")
	}
	if !op2.shutdownCalled {
		t.Error("op2.shutdown was not called")
	}
}

func TestRunner_ExecAndShutdown(t *testing.T) {
	// define timeout
	timeout := 5 * time.Second
	shortTimeout := 100 * time.Millisecond

	// Setup test operations
	op1 := &mockOperation{name: "op1", execDelay: 50 * time.Millisecond}
	op2 := &mockOperation{name: "op2", execDelay: 50 * time.Millisecond}

	execOps := map[string]ExecCallback{
		"op1": op1.exec,
		"op2": op2.exec,
	}

	// Create a cancelable context so we can trigger the shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Run with a generous timeout
	wait, err := Runner(ctx, timeout, execOps)
	if err != nil {
		t.Fatalf("Runner failed with error: %v", err)
	}

	// Verify that exec was called for both operations
	if !op1.execCalled {
		t.Error("op1.exec was not called")
	}
	if !op2.execCalled {
		t.Error("op2.exec was not called")
	}

	// Trigger the shutdown after a short delay
	go func() {
		time.Sleep(shortTimeout)
		// This simulates a shutdown signal
		cancel()
	}()

	// Wait for shutdown to complete
	select {
	case <-wait:
		// Success
	case <-time.After(timeout):
		t.Fatal("Shutdown did not complete within expected time")
	}

	// Verify shutdown was called for both
	if !op1.shutdownCalled {
		t.Error("op1.shutdown was not called")
	}
	if !op2.shutdownCalled {
		t.Error("op2.shutdown was not called")
	}
}
