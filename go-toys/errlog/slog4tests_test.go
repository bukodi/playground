package errlog

import (
	"log/slog"
	"testing"
)

func TestTestSLogHandler_Enabled(t *testing.T) {
	slog.Info("before capture")
	logRecords := CaptureSLog(t, func() {
		slog.Info("inside capture")
	})
	slog.Info("after capture")

	for _, r := range logRecords {
		t.Logf("captured record: %v", r)
	}
}
