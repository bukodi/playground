package errlog

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestConsoleHandler(t *testing.T) {
	var buf bytes.Buffer
	opts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}
	handler := NewConsoleHandler(&buf, opts)
	logger := slog.New(handler)

	// Test different log levels
	logger.Debug("Debug message", "key1", "value1")
	logger.Info("Info message", "key2", "value2")
	logger.Warn("Warning message", "key3", "value3")
	logger.Error("Error message", "key4", "value4")

	output := buf.String()
	t.Logf("Log output:\n%s", output)

	// Check that all messages are present
	if !strings.Contains(output, "Debug message") {
		t.Error("Debug message not found in output")
	}
	if !strings.Contains(output, "Info message") {
		t.Error("Info message not found in output")
	}
	if !strings.Contains(output, "Warning message") {
		t.Error("Warning message not found in output")
	}
	if !strings.Contains(output, "Error message") {
		t.Error("Error message not found in output")
	}

	// Check that attributes are present
	if !strings.Contains(output, "key1=value1") {
		t.Error("key1=value1 not found in output")
	}
	if !strings.Contains(output, "key2=value2") {
		t.Error("key2=value2 not found in output")
	}
	if !strings.Contains(output, "key3=value3") {
		t.Error("key3=value3 not found in output")
	}
	if !strings.Contains(output, "key4=value4") {
		t.Error("key4=value4 not found in output")
	}

	// Check that source information is present
	if !strings.Contains(output, "source=") {
		t.Error("Source information not found in output")
	}
}

func TestConsoleHandlerWithAttrs(t *testing.T) {
	var buf bytes.Buffer
	handler := NewConsoleHandler(&buf, nil)

	// Add some attributes
	handler = handler.WithAttrs([]slog.Attr{
		slog.String("app", "test"),
		slog.Int("version", 1),
	}).(*ConsoleHandler)

	logger := slog.New(handler)
	logger.Info("Message with attributes")

	output := buf.String()
	t.Logf("Log output:\n%s", output)

	// Check that the attributes are present
	if !strings.Contains(output, "app=test") {
		t.Error("app=test not found in output")
	}
	if !strings.Contains(output, "version=1") {
		t.Error("version=1 not found in output")
	}
}

func TestConsoleHandlerWithGroup(t *testing.T) {
	var buf bytes.Buffer
	handler := NewConsoleHandler(&buf, nil)

	// Add a group
	handler = handler.WithGroup("group1").(*ConsoleHandler)

	logger := slog.New(handler)
	logger.Info("Message with group", "key", "value")

	output := buf.String()
	t.Logf("Log output:\n%s", output)

	// Check that the group is present
	if !strings.Contains(output, "group1.key=value") {
		t.Error("group1.key=value not found in output")
	}
}

func TestConsoleHandlerEnabled(t *testing.T) {
	// Test with default level (Info)
	handler := NewConsoleHandler(nil, nil)

	if handler.Enabled(context.Background(), slog.LevelDebug) {
		t.Error("Handler should not be enabled for Debug level with default options")
	}
	if !handler.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("Handler should be enabled for Info level with default options")
	}
	if !handler.Enabled(context.Background(), slog.LevelWarn) {
		t.Error("Handler should be enabled for Warn level with default options")
	}
	if !handler.Enabled(context.Background(), slog.LevelError) {
		t.Error("Handler should be enabled for Error level with default options")
	}

	// Test with custom level
	opts := &slog.HandlerOptions{
		Level: slog.LevelWarn,
	}
	handler = NewConsoleHandler(nil, opts)

	if handler.Enabled(context.Background(), slog.LevelDebug) {
		t.Error("Handler should not be enabled for Debug level with Warn level option")
	}
	if handler.Enabled(context.Background(), slog.LevelInfo) {
		t.Error("Handler should not be enabled for Info level with Warn level option")
	}
	if !handler.Enabled(context.Background(), slog.LevelWarn) {
		t.Error("Handler should be enabled for Warn level with Warn level option")
	}
	if !handler.Enabled(context.Background(), slog.LevelError) {
		t.Error("Handler should be enabled for Error level with Warn level option")
	}
}

func TestConsoleHandlerRecord(t *testing.T) {
	var buf bytes.Buffer
	handler := NewConsoleHandler(&buf, nil)

	// Create a record directly
	record := slog.Record{
		Time:    time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC),
		Message: "Direct record",
		Level:   slog.LevelInfo,
	}
	record.AddAttrs(slog.String("key", "value"))

	// Handle the record
	err := handler.Handle(context.Background(), record)
	if err != nil {
		t.Fatalf("Handler.Handle returned error: %v", err)
	}

	output := buf.String()
	t.Logf("Log output:\n%s", output)

	// Check that the record was logged correctly
	if !strings.Contains(output, "2023-01-01T12:00:00.000Z") {
		t.Error("Time not found in output")
	}
	if !strings.Contains(output, "INFO") {
		t.Error("Level not found in output")
	}
	if !strings.Contains(output, "Direct record") {
		t.Error("Message not found in output")
	}
	if !strings.Contains(output, "key=value") {
		t.Error("key=value not found in output")
	}
}
