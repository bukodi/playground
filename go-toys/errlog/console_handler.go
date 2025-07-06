package errlog

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"sync"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[37m"
	colorWhite  = "\033[97m"
)

// ConsoleHandler is a slog.Handler that writes log records to the console with colorized output.
type ConsoleHandler struct {
	w           io.Writer
	opts        slog.HandlerOptions
	mu          *sync.Mutex
	attrs       []slog.Attr
	groups      []string
	colorize    bool
	levelColors map[slog.Level]string
}

// Enabled reports whether the handler handles records at the given level.
func (c *ConsoleHandler) Enabled(ctx context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if c.opts.Level != nil {
		minLevel = c.opts.Level.Level()
	}
	return level >= minLevel
}

// Handle formats and writes a log record to the console with colorized output.
func (c *ConsoleHandler) Handle(ctx context.Context, record slog.Record) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Format the log record
	var buf []byte

	// Add time if not zero
	if !record.Time.IsZero() {
		timeStr := record.Time.Format("2006-01-02T15:04:05.000Z07:00")
		buf = append(buf, timeStr...)
		buf = append(buf, ' ')
	}

	// Add level with color
	levelColor := c.getLevelColor(record.Level)
	if c.colorize {
		buf = append(buf, levelColor...)
	}
	buf = append(buf, record.Level.String()...)
	if c.colorize {
		buf = append(buf, colorReset...)
	}
	buf = append(buf, ' ')

	// Add message
	buf = append(buf, record.Message...)

	// Add source if enabled
	if c.opts.AddSource && record.PC != 0 {
		frames := runtime.CallersFrames([]uintptr{record.PC})
		frame, _ := frames.Next()
		if frame.File != "" {
			buf = append(buf, ' ')
			buf = append(buf, fmt.Sprintf("source=%s:%d", frame.File, frame.Line)...)
		}
	}

	// Add attributes
	for _, attr := range c.attrs {
		buf = append(buf, ' ')
		buf = append(buf, attr.Key...)
		buf = append(buf, '=')
		buf = append(buf, fmt.Sprintf("%v", attr.Value.Any())...)
	}

	// Add attributes from the record
	record.Attrs(func(attr slog.Attr) bool {
		buf = append(buf, ' ')

		// Add group prefix to key if there are groups
		if len(c.groups) > 0 {
			for i, group := range c.groups {
				buf = append(buf, group...)
				if i < len(c.groups)-1 || attr.Key != "" {
					buf = append(buf, '.')
				}
			}
		}

		buf = append(buf, attr.Key...)
		buf = append(buf, '=')
		buf = append(buf, fmt.Sprintf("%v", attr.Value.Any())...)
		return true
	})

	buf = append(buf, '\n')

	_, err := c.w.Write(buf)
	return err
}

// WithAttrs returns a new ConsoleHandler with the given attributes added.
func (c *ConsoleHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return c
	}

	newHandler := c.clone()
	newHandler.attrs = append(sliceClone(c.attrs), attrs...)
	return newHandler
}

// WithGroup returns a new ConsoleHandler with the given group added.
func (c *ConsoleHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return c
	}

	newHandler := c.clone()
	newHandler.groups = append(sliceClone(c.groups), name)
	return newHandler
}

// Type guard
var _ slog.Handler = &ConsoleHandler{}

// NewConsoleHandler creates a new ConsoleHandler that writes to w.
func NewConsoleHandler(w io.Writer, opts *slog.HandlerOptions) *ConsoleHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	return &ConsoleHandler{
		w:        w,
		opts:     *opts,
		mu:       &sync.Mutex{},
		colorize: true,
		levelColors: map[slog.Level]string{
			slog.LevelDebug: colorGray,
			slog.LevelInfo:  colorGreen,
			slog.LevelWarn:  colorYellow,
			slog.LevelError: colorRed,
		},
	}
}

// clone creates a copy of the handler for use by WithAttrs and WithGroup.
func (c *ConsoleHandler) clone() *ConsoleHandler {
	return &ConsoleHandler{
		w:           c.w,
		opts:        c.opts,
		mu:          c.mu, // mutex is shared among all clones
		attrs:       sliceClone(c.attrs),
		groups:      sliceClone(c.groups),
		colorize:    c.colorize,
		levelColors: c.levelColors,
	}
}

// getLevelColor returns the color for the given level.
func (c *ConsoleHandler) getLevelColor(level slog.Level) string {
	if color, ok := c.levelColors[level]; ok {
		return color
	}

	// Default colors based on level
	switch {
	case level < slog.LevelInfo:
		return colorGray
	case level < slog.LevelWarn:
		return colorGreen
	case level < slog.LevelError:
		return colorYellow
	default:
		return colorRed
	}
}

// sliceClone returns a copy of the slice.
func sliceClone[T any](s []T) []T {
	if s == nil {
		return nil
	}
	return append(make([]T, 0, len(s)), s...)
}
