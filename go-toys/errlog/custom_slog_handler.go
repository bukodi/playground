package errlog

import (
	"context"
	"log/slog"
)

type MyHandler struct {
}

func (m MyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	//TODO implement me
	panic("implement me")
}

func (m MyHandler) Handle(ctx context.Context, record slog.Record) error {
	//TODO implement me
	panic("implement me")
}

func (m MyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	//TODO implement me
	panic("implement me")
}

func (m MyHandler) WithGroup(name string) slog.Handler {
	//TODO implement me
	panic("implement me")
}

var _ slog.Handler = &MyHandler{}
