package tbluge

import (
	"github.com/blugelabs/bluge"
	"github.com/blugelabs/bluge/index"
	"log"
	"testing"
)

func TestBluge(t *testing.T) {
	config := index.InMemoryOnlyConfig()
	writer, err := index.OpenWriter(config)
	if err != nil {
		log.Fatalf("error opening writer: %v", err)
	}
	defer writer.Close()

	doc := bluge.NewDocument("example").
		AddField(bluge.NewTextField("name", "bluge"))

	batch := index.NewBatch()
	batch.Insert(doc)
	err = writer.Batch(batch)
	if err != nil {
		log.Fatalf("error updating document: %v", err)
	}
}
