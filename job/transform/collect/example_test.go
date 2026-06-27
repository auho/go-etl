package collect

import (
	"github.com/auho/go-etl/v3/job/extract"
)

func ExampleNewKeysAll() {
	// NewKeysAll creates a collector that collects values from all specified keys
	collector := NewKeysAll([]string{"name", "email"})

	// Collector is typically used with filter.NewFilterPredicate, transform.NewInsert, etc.
	// e.g., filter.NewFilterPredicate(collector, extractor)
	var _ extract.Extractor
	_ = collector
}

func ExampleNewKeysAny() {
	// NewKeysAny creates a collector that tries each key in order
	// and returns on the first matching result
	collector := NewKeysAny([]string{"phone", "mobile"})

	var _ extract.Extractor
	_ = collector
}
