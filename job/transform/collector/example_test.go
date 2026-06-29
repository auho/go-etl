package collector_test

import (
	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collector"
	"github.com/auho/go-etl/v3/job/transform/collector/keys"
	"github.com/auho/go-etl/v3/job/transform/collector/mode"
)

func ExampleNewCollector() {
	// NewCollector creates a collector with a Keys source, All mode, and an extractor.
	collector := collector.NewCollector(keys.New([]string{"name", "email"}), mode.NewAll(), nil /* extractor */)

	var _ extract.Extractor
	_ = collector
}