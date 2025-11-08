package tctx

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"testing"
	"testing/synctest"
	"time"
)

var testTask = func(ctx context.Context, name string, dur time.Duration) error {
	slog.Info(name + " started")
	select {
	case <-ctx.Done():
		slog.Info(fmt.Sprintf("%s canceled with: %+v", name, ctx.Err()))
		return ctx.Err()
	case <-time.After(dur):
		slog.Info(name + " finished")
		return nil
	}
}

func TestExecuteSimple(t *testing.T) {
	timeUnit := time.Second
	ste := &SingleThreadExecutor{}

	synctest.Test(t, func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Go(func() {
			err := ste.Execute(t.Context(), 20*timeUnit, func(ctx context.Context) error {
				return testTask(ctx, "task1", 10*timeUnit)
			})
			t.Logf("err = %v", err)
		})
		wg.Go(func() {
			time.Sleep(2 * timeUnit)
			err := ste.Execute(t.Context(), 20*timeUnit, func(ctx context.Context) error {
				return testTask(ctx, "task2", 5*timeUnit)
			})
			t.Logf("err = %v", err)
		})
		wg.Wait()
	})
}
