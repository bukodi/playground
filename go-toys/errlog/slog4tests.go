package errlog

import (
	"bytes"
	"context"
	"log/slog"
	"sync"
	"testing"
)

type TestSLogHandler struct {
	t                *testing.T
	buffer           *bytes.Buffer
	infoTextHandler  slog.Handler
	errorTextHandler slog.Handler
	originalHandler  slog.Handler
	mu               sync.Mutex
	logRecords       []slog.Record
}

type writerFunc struct {
	fn func(p []byte) (n int, err error)
}

func (wf writerFunc) Write(p []byte) (n int, err error) {
	return wf.fn(p)
}

func NewHandler(t *testing.T) *TestSLogHandler {
	tlh := &TestSLogHandler{t: t}

	tlh.buffer = &bytes.Buffer{}
	tlh.infoTextHandler = slog.NewTextHandler(writerFunc{
		fn: func(p []byte) (n int, err error) {
			t.Logf("out: %s", p)
			return len(p), nil
		},
	}, &slog.HandlerOptions{
		AddSource: true,
	})
	tlh.errorTextHandler = slog.NewTextHandler(writerFunc{
		fn: func(p []byte) (n int, err error) {
			t.Errorf("err: %s", p)
			return len(p), nil
		},
	}, &slog.HandlerOptions{
		AddSource: true,
	})
	tlh.originalHandler = slog.Default().Handler()
	slog.SetDefault(slog.New(tlh))
	return tlh
}

func (tlh *TestSLogHandler) Close() []slog.Record {
	tlh.mu.Lock()
	defer tlh.mu.Unlock()

	slog.SetDefault(slog.New(tlh.originalHandler))
	tlh.originalHandler = nil
	return tlh.logRecords
}

func CaptureSLog(t *testing.T, fn func()) []slog.Record {
	tlh := NewHandler(t)
	fn()
	return tlh.Close()
}

func (tlh *TestSLogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (tlh *TestSLogHandler) Handle(ctx context.Context, record slog.Record) error {
	tlh.mu.Lock()
	defer tlh.mu.Unlock()
	tlh.logRecords = append(tlh.logRecords, record)

	if record.Level >= slog.LevelError {
		return tlh.errorTextHandler.Handle(ctx, record)
	} else {
		return tlh.infoTextHandler.Handle(ctx, record)
	}
}

func (tlh *TestSLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	panic("implement me")
}

func (tlh *TestSLogHandler) WithGroup(name string) slog.Handler {
	panic("implement me")
}

var _ slog.Handler = &TestSLogHandler{}
