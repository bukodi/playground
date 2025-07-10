package totter

import "testing"

import (
	"context"
	"time"

	"github.com/maypok86/otter/v2"
	"github.com/maypok86/otter/v2/stats"
)

func TestOtter(t *testing.T) {
	ctx := t.Context()

	// Create statistics counter to track cache operations
	counter := stats.NewCounter()

	// Configure cache with:
	// - Capacity: 10,000 entries
	// - 1 second expiration after last access
	// - 500ms refresh interval after writes
	// - Stats collection enabled
	cache := otter.Must(&otter.Options[string, string]{
		MaximumSize:       10_000,
		ExpiryCalculator:  otter.ExpiryAccessing[string, string](time.Second),           // Reset timer on reads/writes
		RefreshCalculator: otter.RefreshWriting[string, string](500 * time.Millisecond), // Refresh after writes
		StatsRecorder:     counter,                                                      // Attach stats collector
	})

	// Phase 1: Test basic expiration
	// -----------------------------
	cache.Set("key", "value") // Add initial value

	// Wait for expiration (1 second)
	time.Sleep(time.Second)

	// Verify entry expired
	if _, ok := cache.GetIfPresent("key"); ok {
		t.Fatalf("key shouldn't be found") // Should be expired
	}

	// Phase 2: Test cache stampede protection
	// --------------------------------------
	loader := func(ctx context.Context, key string) (string, error) {
		time.Sleep(200 * time.Millisecond) // Simulate slow load
		return "value1", nil               // Return new value
	}

	// Concurrent Gets would deduplicate loader calls
	value, err := cache.Get(ctx, "key", otter.LoaderFunc[string, string](loader))
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if value != "value1" {
		t.Fatalf("incorrect value") // Should get newly loaded value
	}

	// Phase 3: Test background refresh
	// --------------------------------
	time.Sleep(500 * time.Millisecond) // Wait until refresh needed

	// New loader that returns updated value
	loader = func(ctx context.Context, key string) (string, error) {
		time.Sleep(100 * time.Millisecond) // Simulate refresh
		return "value2", nil               // Return refreshed value
	}

	// This triggers async refresh but returns current value
	value, err = cache.Get(ctx, "key", otter.LoaderFunc[string, string](loader))
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if value != "value1" { // Should get old value while refreshing
		t.Fatalf("loader shouldn't be called during Get")
	}

	// Wait for refresh to complete
	time.Sleep(110 * time.Millisecond)

	// Verify refreshed value
	v, ok := cache.GetIfPresent("key")
	if !ok {
		t.Fatalf("key should be found") // Should still be cached
	}
	if v != "value2" { // Should now have refreshed value
		t.Fatalf("refresh should be completed")
	}
}
