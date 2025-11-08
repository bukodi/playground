package tctx

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// SingleThreadExecutor runs at most one task at a time.
// Submitting a new task cancels the one currently running.
// It is safe for concurrent use by multiple goroutines.
//
// Notes:
//   - Tasks must respect ctx.Done() to be preemptible.
//   - Execute runs the task synchronously (call returns when task returns).
//     If you want fire-and-forget behavior, call Execute from your own goroutine.
//
// Typical usage:
//   ste := &SingleThreadExecutor{}
//   go ste.Execute(ctx, 10*time.Second, func(ctx context.Context) error { /* ... */ return nil })
//   // later submit another task; the previous one will be canceled
//   go ste.Execute(ctx, 5*time.Second, func(ctx context.Context) error { /* ... */ return nil })
//
// If you need to explicitly stop whatever is running without submitting a replacement,
// call ste.Stop().

type SingleThreadExecutor struct {
	mu       sync.Mutex
	ctx      context.Context
	cancelFn context.CancelCauseFunc
}

// Stop cancels the currently running task, if any.
func (e *SingleThreadExecutor) Stop() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.stopLocked(errors.New("executor stopped"))
}

// stopLocked cancels the current task context. Caller must hold e.mu.
func (e *SingleThreadExecutor) stopLocked(cause error) {
	if e.ctx == nil {
		return
	}
	// Cancel current task with a cause.
	e.cancelFn(cause)
	// Clear references so a future Execute sees no active ctx.
	e.ctx = nil
	e.cancelFn = nil
}

// Execute cancels any running task, then starts and runs the provided task synchronously.
// - parentCtx: the parent for the task context
// - timeout: per-task timeout
// - task: the function to run; must honor ctx.Done() for prompt cancellation
// Returns the task's error, or context errors (canceled/timeout) if preempted or timed out.
func (e *SingleThreadExecutor) Execute(parentCtx context.Context, timeout time.Duration, task func(ctx context.Context) error) error {
	// Preempt any currently running task and prepare a fresh context for this one.
	e.mu.Lock()
	// Cancel whatever is running.
	e.stopLocked(fmt.Errorf("preempted by new task"))

	// Create a timeout context, then a cancelable child we can cancel with a cause.
	// If timeout <= 0, use a no-timeout child of parent.
	var base context.Context
	var cancelTimeout context.CancelFunc
	if timeout > 0 {
		base, cancelTimeout = context.WithTimeoutCause(parentCtx, timeout, errors.New("task timeout"))
	} else {
		base, cancelTimeout = context.WithCancel(parentCtx)
	}
	ctx, cancelWithCause := context.WithCancelCause(base)
	// Store current task context and cancel function.
	e.ctx = ctx
	e.cancelFn = cancelWithCause
	e.mu.Unlock()

	// Ensure the timeout context is cleaned up when task exits.
	defer cancelTimeout()

	// Run the task synchronously. It should return promptly if ctx is canceled.
	err := task(ctx)

	// If this task is still the current one, clear it and make sure it's canceled
	// to unblock any waiters listening on ctx.Done().
	e.mu.Lock()
	if e.ctx == ctx {
		// Clear and cancel with the task's own error as the cause if not already canceled.
		// If err is nil and ctx is Done(), keep the existing cause.
		cause := err
		if cause == nil {
			// If the context completed due to timeout/cancel, propagate that as cause.
			select {
			case <-ctx.Done():
				cause = context.Cause(ctx)
			default:
				// still nil
			}
		}
		e.stopLocked(cause)
	}
	e.mu.Unlock()

	return err
}
