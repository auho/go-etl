package collector_test

import (
	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collector"
)

func ExampleNewCollector() {
	// NewCollector creates a collector with a Keys source, All mode, and an extractor.
	collector := collector.NewKeysAll([]string{"name", "email"}, nil /* extractor */)

	var _ extract.Extractor
	_ = collector
}