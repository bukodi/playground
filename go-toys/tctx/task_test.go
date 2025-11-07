package tctx

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

type singleThreadExecutor struct {
	mu       sync.Mutex
	ctx      context.Context
	cancelFn context.CancelFunc
}

var worker = singleThreadExecutor{}

func (ste *singleThreadExecutor) Execute(parentCtx context.Context, f func(ctx context.Context), timeout time.Duration) {
	ste.mu.Lock()
	ste.ctx, ste.cancelFn = context.WithTimeoutCause(parentCtx, timeout, fmt.Errorf("timeout"))
	context.AfterFunc(ste.ctx, func() {
		fmt.Printf("after func")
	})
	ste.mu.Unlock()
	f(ste.ctx)
	ste.mu.Lock()
	ste.cancelFn()
	ste.mu.Unlock()
}

func (ste *singleThreadExecutor) Stop() {
	ste.mu.Lock()
	ste.cancelFn()
	ste.mu.Unlock()
}

func TestExecuteTask(t *testing.T) {
	worker.Execute(t.Context(), func(ctx context.Context) {
		t.Logf("Task started")
		select {
		case <-ctx.Done():
			t.Logf("Task canceled")
			return
		case <-time.Tick(2 * time.Second):
			t.Logf("Task finished")
		}
	}, 1*time.Second)
	worker.Stop()
}
