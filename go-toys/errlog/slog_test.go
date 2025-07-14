package errlog

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"os"
	"testing"
	"testing/slogtest"
)

type forwardingHandler struct {
	parentHandler func() slog.Handler
}

func (m forwardingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return m.parentHandler().Enabled(ctx, level)
}

func (m forwardingHandler) Handle(ctx context.Context, record slog.Record) error {
	return m.parentHandler().Handle(ctx, record)
}

func (m forwardingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return m.parentHandler().WithAttrs(attrs)
}

func (m forwardingHandler) WithGroup(name string) slog.Handler {
	return m.parentHandler().WithGroup(name)
}

var _ slog.Handler = &forwardingHandler{}

var pkgLogger = func() *slog.Logger {
	fh := forwardingHandler{parentHandler: func() slog.Handler { return slog.Default().Handler() }}
	return slog.New(fh)
}()

//var pkgLogger = NewPkgLogger(nil)

func foo(ctx context.Context) {
	slog.Info("info from foo")
	slog.InfoContext(ctx, "info from foo with ctx")
	if err := bar(ctx); err != nil {
		slog.Error("bar failed", "err", err)
	}

	slog.Debug("big data", "dump", func() slog.Value {
		return slog.StringValue("dynamic data")
	})
}

func bar(ctx context.Context) error {
	slog.Info("info from the bar")
	slog.InfoContext(ctx, "info from bar with ctx")
	slog.Debug("debug from bar")
	slog.DebugContext(ctx, "debug from bar with ctx")
	return errors.New("an error occurred in the bar")
}

func TestSlogPlain(t *testing.T) {
	ctx := context.TODO()
	foo(ctx)
}

func TestSlogLevelerVar(t *testing.T) {
	ctx := context.TODO()
	foo(ctx)

	// Change the logger
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   true,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	})
	defLogger := slog.New(jsonHandler)
	slog.SetDefault(defLogger)

	foo(ctx) //pkgLogger = NewPkgLogger(nil)

	dynamicLevel := &slog.LevelVar{}

	// Change again
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		Level:       dynamicLevel,
		ReplaceAttr: nil,
	})))

	foo(ctx) //pkgLogger = NewPkgLogger(nil)

	slog.Info("top level info")
	slog.Debug("top level debug")
	dynamicLevel.Set(slog.LevelDebug)
	if result, err := dynamicLevel.AppendText([]byte("cica")); err != nil {
		t.Fatalf("error: %s", err)
	} else {
		t.Logf("result: %s", result)
	}

	slog.Info("top level info after dynamicLevel = slog.LevelDebug")
	slog.Debug("top level debug after dynamicLevel = slog.LevelDebug")

}

func TestSlogHandler(t *testing.T) {
	var buf bytes.Buffer
	h := slog.NewJSONHandler(&buf, nil)

	results := func() []map[string]any {
		var ms []map[string]any
		for _, line := range bytes.Split(buf.Bytes(), []byte{'\n'}) {
			if len(line) == 0 {
				continue
			}
			var m map[string]any
			if err := json.Unmarshal(line, &m); err != nil {
				panic(err) // In a real test, use t.Fatal.
			}
			ms = append(ms, m)
		}
		return ms
	}
	err := slogtest.TestHandler(h, results)
	if err != nil {
		log.Fatal(err)
	}
}

func TestSlogHandlerWithAttrs(t *testing.T) {
	//slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})))
	slog.Info("without attrs")
	slog.Info("with one attr", slog.String("key", "value"))
	slog.Info("with empty attr", slog.Attr{})
	slog.Info("with empty and one attr", slog.Attr{}, slog.String("key", "value"))

	ctx := t.Context()
	slog.InfoContext(ctx, "with context and one attr", slog.String("key", "value"))

	slog.Info("with groupAttr", slog.Group("group1", slog.String("groupKey1", "value 1"), slog.String("groupKey1", "value 1")))

	slog.Info("with groupAttr without group name", slog.Group("", slog.String("groupKey1", "value 1"), slog.String("groupKey1", "value 1")))

	slog.Info("with a duplicate key", slog.String("key", "value1"), slog.String("key", "value2"))

}
