package errlog

import (
	"context"
	"log/slog"
	"strings"
	"testing"
	"time"
)

// Helper function to find an attribute in a slog.Record
func findAttr(record slog.Record, key string) (string, bool) {
	var value string
	var found bool

	// For group1.key format, we need to split the key
	parts := strings.Split(key, ".")
	if len(parts) > 1 {
		// This is a grouped key, we need to check if the attribute has the group prefix
		groupPrefix := parts[0]
		attrKey := parts[1]

		record.Attrs(func(attr slog.Attr) bool {
			if strings.HasPrefix(attr.Key, groupPrefix+".") && strings.TrimPrefix(attr.Key, groupPrefix+".") == attrKey {
				value = attr.Value.String()
				found = true
				return false
			}
			return true
		})
	} else {
		// Regular key
		record.Attrs(func(attr slog.Attr) bool {
			if attr.Key == key {
				value = attr.Value.String()
				found = true
				return false
			}
			return true
		})
	}

	return value, found
}

func TestInMemoryHandler(t *testing.T) {
	opts := &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}
	handler := NewInMemoryHandler(10, opts)
	logger := slog.New(handler)

	// Test different log levels
	logger.Debug("Debug message", "key1", "value1")
	logger.Info("Info message", "key2", "value2")
	logger.Warn("Warning message", "key3", "value3")
	logger.Error("Error message", "key4", "value4")

	records := handler.GetRecords()
	t.Logf("Records: %v", records)

	// Check that all messages are present
	if len(records) != 4 {
		t.Errorf("Expected 4 records, got %d", len(records))
	}

	// Check the messages and levels
	expectedMessages := []string{
		"Debug message",
		"Info message",
		"Warning message",
		"Error message",
	}
	expectedLevels := []slog.Level{
		slog.LevelDebug,
		slog.LevelInfo,
		slog.LevelWarn,
		slog.LevelError,
	}
	expectedKeys := []string{
		"key1",
		"key2",
		"key3",
		"key4",
	}
	expectedValues := []string{
		"value1",
		"value2",
		"value3",
		"value4",
	}

	for i, record := range records {
		if record.Message != expectedMessages[i] {
			t.Errorf("Expected message %q, got %q", expectedMessages[i], record.Message)
		}
		if record.Level != expectedLevels[i] {
			t.Errorf("Expected level %v, got %v", expectedLevels[i], record.Level)
		}
		if value, ok := findAttr(record, expectedKeys[i]); !ok || value != expectedValues[i] {
			t.Errorf("Expected attribute %s=%s, got %s", expectedKeys[i], expectedValues[i], value)
		}
	}
}

func TestInMemoryHandlerWithAttrs(t *testing.T) {
	handler := NewInMemoryHandler(10, nil)

	// Add some attributes
	handler = handler.WithAttrs([]slog.Attr{
		slog.String("app", "test"),
		slog.Int("version", 1),
	}).(*InMemoryHandler)

	logger := slog.New(handler)
	logger.Info("Message with attributes")

	records := handler.GetRecords()
	t.Logf("Records: %v", records)

	if len(records) != 1 {
		t.Errorf("Expected 1 record, got %d", len(records))
	}

	// Check that the attributes are present
	record := records[0]
	if value, ok := findAttr(record, "app"); !ok || value != "test" {
		t.Errorf("Expected attribute app=test, got %s", value)
	}
	if value, ok := findAttr(record, "version"); !ok || value != "1" {
		t.Errorf("Expected attribute version=1, got %s", value)
	}
}

func TestInMemoryHandlerWithGroup(t *testing.T) {
	handler := NewInMemoryHandler(10, nil)

	// Add a group
	handler = handler.WithGroup("group1").(*InMemoryHandler)

	logger := slog.New(handler)
	logger.Info("Message with group", "key", "value")

	records := handler.GetRecords()
	t.Logf("Records: %v", records)

	if len(records) != 1 {
		t.Errorf("Expected 1 record, got %d", len(records))
	}

	// Check that the group is present
	record := records[0]
	if value, ok := findAttr(record, "group1.key"); !ok || value != "value" {
		t.Errorf("Expected attribute group1.key=value, got %s", value)
	}
}

func TestInMemoryHandlerEnabled(t *testing.T) {
	// Test with default level (Info)
	handler := NewInMemoryHandler(10, nil)

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
	handler = NewInMemoryHandler(10, opts)

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

func TestInMemoryHandlerRecord(t *testing.T) {
	handler := NewInMemoryHandler(10, nil)

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

	records := handler.GetRecords()
	t.Logf("Records: %v", records)

	if len(records) != 1 {
		t.Errorf("Expected 1 record, got %d", len(records))
	}

	// Check that the record was stored correctly
	storedRecord := records[0]
	expectedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	if !storedRecord.Time.Equal(expectedTime) {
		t.Errorf("Expected time %v, got %v", expectedTime, storedRecord.Time)
	}
	if storedRecord.Level != slog.LevelInfo {
		t.Errorf("Expected level INFO, got %v", storedRecord.Level)
	}
	if storedRecord.Message != "Direct record" {
		t.Errorf("Expected message 'Direct record', got %s", storedRecord.Message)
	}
	if value, ok := findAttr(storedRecord, "key"); !ok || value != "value" {
		t.Errorf("Expected attribute key=value, got %s", value)
	}
}

func TestInMemoryHandlerMaxSize(t *testing.T) {
	// Create a handler with a small max size
	handler := NewInMemoryHandler(3, nil)
	logger := slog.New(handler)

	// Log more messages than the max size
	logger.Info("Message 1")
	logger.Info("Message 2")
	logger.Info("Message 3")
	logger.Info("Message 4")
	logger.Info("Message 5")

	records := handler.GetRecords()
	t.Logf("Records: %v", records)

	// Check that only the most recent messages are kept
	if len(records) != 3 {
		t.Errorf("Expected 3 records, got %d", len(records))
	}

	expectedMessages := []string{
		"Message 3",
		"Message 4",
		"Message 5",
	}

	for i, record := range records {
		if record.Message != expectedMessages[i] {
			t.Errorf("Expected message %q, got %q", expectedMessages[i], record.Message)
		}
	}
}
