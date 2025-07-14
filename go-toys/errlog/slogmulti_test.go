package errlog

import (
	"bytes"
	"context"
	slogmulti "github.com/samber/slog-multi"
	"gopkg.in/natefinch/lumberjack.v2"
	"log/slog"
	"path/filepath"
	"sync"
	"testing"
)

type inMemoryCollector struct {
	records []slog.Record
	mu      *sync.Mutex
	maxSize int
}

func (imc *inMemoryCollector) handleRecord(ctx context.Context, groups []string, attrs []slog.Attr, record slog.Record) error {
	imc.mu.Lock()
	defer imc.mu.Unlock()

	// If buffer is full, remove the oldest record
	if imc.maxSize > 0 && len(imc.records) >= imc.maxSize {
		imc.records = imc.records[1:]
	}
	imc.records = append(imc.records, record)
	return nil
}

func TestSlogMulti(t *testing.T) {
	imc := &inMemoryCollector{
		records: make([]slog.Record, 0, 100),
		mu:      &sync.Mutex{},
		maxSize: 100,
	}

	var textBuff bytes.Buffer
	var jsonBuff bytes.Buffer

	var dynamicLevel slog.LevelVar

	fanOutHandler := slogmulti.Fanout(
		slog.NewTextHandler(&textBuff, &slog.HandlerOptions{
			AddSource: true,
			Level:     &dynamicLevel,
		}),
		slog.NewJSONHandler(&jsonBuff, &slog.HandlerOptions{
			AddSource: true,
			Level:     &dynamicLevel,
		}),
		slog.NewTextHandler(&lumberjack.Logger{
			Filename:   filepath.Join(t.TempDir(), "test.log"),
			MaxSize:    1, // megabytes
			MaxBackups: 5,
			MaxAge:     3, //days
			Compress:   false,
		}, &slog.HandlerOptions{
			AddSource: true,
			Level:     &dynamicLevel,
		}),
		slogmulti.NewHandleInlineHandler(imc.handleRecord),
	)

	defaultHandler := slogmulti.Pipe(slogmulti.NewHandleInlineMiddleware(fnProcessAttrs)).Handler(fanOutHandler)

	slog.SetDefault(slog.New(defaultHandler))
	slog.Info("first info", slog.String("attrKey", "attrValue"))
	slog.Info("second info", slog.String("attrKey", "attrValue"))
	slog.Info("third info", slog.String("attrKey", "attrValue"))
	dynamicLevel.Set(slog.LevelError)
	slog.Info("second info")
	dynamicLevel.Set(slog.LevelDebug)
	slog.Info("third info")

	t.Logf("\ntext buffer: \n%v", textBuff.String())
	t.Logf("\njson buffer: \n%v", jsonBuff.String())
	t.Logf("\nrecords: \n%v", imc.records)
}

func fnProcessAttrs(ctx context.Context, record slog.Record, next func(context.Context, slog.Record) error) error {
	nr := slog.NewRecord(record.Time, record.Level, record.Message, record.PC)
	attrs := make([]slog.Attr, 0, record.NumAttrs()+1)
	record.Attrs(func(attr slog.Attr) bool {
		attrs = append(attrs, attr)
		return true
	})
	attrs = append(attrs, slog.String("customKey", "customValue"))
	nr.AddAttrs(attrs...)
	return next(ctx, nr)
}
