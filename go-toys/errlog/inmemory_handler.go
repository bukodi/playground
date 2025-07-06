package errlog

import (
	"context"
	"log/slog"
	"sync"
)

// InMemoryHandler is a slog.Handler that stores log records in memory.
type InMemoryHandler struct {
	opts        slog.HandlerOptions
	mu          *sync.Mutex
	attrs       []slog.Attr
	groups      []string
	records     []slog.Record
	maxSize     int
	currentSize int
}

// Enabled reports whether the handler handles records at the given level.
func (i *InMemoryHandler) Enabled(ctx context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if i.opts.Level != nil {
		minLevel = i.opts.Level.Level()
	}
	return level >= minLevel
}

// Handle stores a log record in memory.
func (i *InMemoryHandler) Handle(ctx context.Context, record slog.Record) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	// If buffer is full, remove the oldest record
	if i.currentSize >= i.maxSize {
		i.records = i.records[1:]
		i.currentSize--
	}

	// Clone the record to avoid shared state
	newRecord := record.Clone()

	// Add handler attributes to the record
	for _, attr := range i.attrs {
		newRecord.AddAttrs(attr)
	}

	// If there are groups, we need to create a new record with the group prefixes added to the attributes
	if len(i.groups) > 0 {
		// Create a new record with the same time, level, and message
		groupedRecord := slog.NewRecord(newRecord.Time, newRecord.Level, newRecord.Message, newRecord.PC)

		// Extract all attributes and add them with group prefixes
		var groupedAttrs []slog.Attr
		newRecord.Attrs(func(attr slog.Attr) bool {
			// Create a new key with the group prefix
			key := attr.Key
			for _, group := range i.groups {
				key = group + "." + key
			}
			groupedAttrs = append(groupedAttrs, slog.Attr{Key: key, Value: attr.Value})
			return true
		})

		// Add the grouped attributes to the new record
		groupedRecord.AddAttrs(groupedAttrs...)

		// Use the grouped record instead
		newRecord = groupedRecord
	}

	// Add the record to the buffer
	i.records = append(i.records, newRecord)
	i.currentSize++

	return nil
}

// WithAttrs returns a new InMemoryHandler with the given attributes added.
func (i *InMemoryHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return i
	}

	newHandler := i.clone()
	newHandler.attrs = append(sliceClone(i.attrs), attrs...)
	return newHandler
}

// WithGroup returns a new InMemoryHandler with the given group added.
func (i *InMemoryHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return i
	}

	newHandler := i.clone()
	newHandler.groups = append(sliceClone(i.groups), name)
	return newHandler
}

// GetRecords returns the stored log records.
func (i *InMemoryHandler) GetRecords() []slog.Record {
	i.mu.Lock()
	defer i.mu.Unlock()

	// Return a copy of the records to avoid race conditions
	records := make([]slog.Record, len(i.records))
	for j, record := range i.records {
		records[j] = record.Clone()
	}

	return records
}

// clone creates a copy of the handler for use by WithAttrs and WithGroup.
func (i *InMemoryHandler) clone() *InMemoryHandler {
	return &InMemoryHandler{
		opts:        i.opts,
		mu:          i.mu, // mutex is shared among all clones
		attrs:       sliceClone(i.attrs),
		groups:      sliceClone(i.groups),
		records:     i.records, // records are shared among all clones
		maxSize:     i.maxSize,
		currentSize: i.currentSize,
	}
}

// NewInMemoryHandler creates a new InMemoryHandler with the given maximum size.
func NewInMemoryHandler(maxSize int, opts *slog.HandlerOptions) *InMemoryHandler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}

	if maxSize <= 0 {
		maxSize = 100 // Default size
	}

	return &InMemoryHandler{
		opts:        *opts,
		mu:          &sync.Mutex{},
		records:     make([]slog.Record, 0, maxSize),
		maxSize:     maxSize,
		currentSize: 0,
	}
}

// Type guard
var _ slog.Handler = &InMemoryHandler{}
