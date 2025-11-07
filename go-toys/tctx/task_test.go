package tctx

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

type singleThreadExecutor struct {
	mu       sync.Mutex
	ctx      context.Context
	cancelFn context.CancelCauseFunc
}

func (e *singleThreadExecutor) Execute(parentCtx context.Context, timeout time.Duration, task func(ctx context.Context) error) error {
	e.mu.Lock()
	if e.ctx != nil {
		//
		e.cancelFn(fmt.Errorf("task caceled"))
		select {
		case <-e.ctx.Done():
			slog.DebugContext(e.ctx, "Task canceled after the user restarts")
		case <-time.After(10 * time.Second):
			slog.ErrorContext(e.ctx, "Waiting for task canceled timed out")
		}
		e.ctx = nil
		e.cancelFn = nil
	}
	ctx1, _ := context.WithTimeoutCause(parentCtx, timeout, errors.New("timeout"))
	e.ctx, e.cancelFn = context.WithCancelCause(ctx1)
	e.mu.Unlock()

	ret := task(e.ctx)

	e.mu.Lock()
	e.cancelFn(nil)
	e.ctx = nil
	e.cancelFn = nil
	e.mu.Unlock()
	return ret
}

var testTask = func(ctx context.Context, dur time.Duration) error {
	slog.Info("Task started")
	select {
	case <-ctx.Done():
		slog.Info(fmt.Sprintf("Context canceled with: %+v", ctx.Err()))
		return ctx.Err()
	case <-time.After(dur):
		slog.Info("Task finished")
		return nil
	}
}

func TestExecuteSimple(t *testing.T) {
	timeUnit := time.Second
	ste := &singleThreadExecutor{}

	synctest.Test(t, func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Go(func() {
			err := ste.Execute(t.Context(), 20*timeUnit, func(ctx context.Context) error {
				return testTask(ctx, 10*timeUnit)
			})
			t.Logf("err = %v", err)
		})
		wg.Go(func() {
			time.Sleep(2 * timeUnit)
			err := ste.Execute(t.Context(), 20*timeUnit, func(ctx context.Context) error {
				return testTask(ctx, 5*timeUnit)
			})
			t.Logf("err = %v", err)
		})
		wg.Wait()
	})
}
