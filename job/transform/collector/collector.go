package collector

import (
	"fmt"

	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collector/mode"
	"github.com/auho/go-etl/v3/job/transform/collector/source"
)

// Collector orchestrates a Source, a Mode, and an Extractor.
// All lifecycle and extraction logic is unified here.
type Collector struct {
	source    source.Source
	mode      mode.Mode
	extractor extract.Extractor
}

func NewCollector(s source.Source, m mode.Mode, e extract.Extractor) *Collector {
	return &Collector{source: s, mode: m, extractor: e}
}

// Title returns the combined title of source and extractor.
func (c *Collector) Title() string {
	return fmt.Sprintf("%s | %s", c.source.Title(), c.extractor.Title())
}

// Fields returns the source keys (input fields for source data row selection).
func (c *Collector) Fields() []string {
	return c.source.Keys()
}

// Keys returns the extractor export keys (output fields).
func (c *Collector) Keys() []string {
	export := c.extractor.NewExport()
	if export == nil {
		return nil
	}

	return export.Keys()
}

// DefaultValues returns the extractor export default values.
func (c *Collector) DefaultValues() map[string]any {
	return c.extractor.NewExport().DefaultValues()
}

func (c *Collector) Prepare() error {
	if err := c.source.Prepare(); err != nil {
		return fmt.Errorf("source.Prepare: %w", err)
	}

	if err := c.mode.Prepare(); err != nil {
		return fmt.Errorf("mode.Prepare: %w", err)
	}

	if err := c.extractor.Prepare(); err != nil {
		return fmt.Errorf("extractor.Prepare: %w", err)
	}

	return nil
}

func (c *Collector) Close() error {
	return c.extractor.Close()
}

// Extract fetches contents from the source for the given item, then applies
// the mode to drive the extractor and returns the result.
func (c *Collector) Extract(item map[string]any) (extract.Result, error) {
	keys, keysValue, err := c.source.Contents(item)
	if err != nil {
		return extract.Result{}, fmt.Errorf("source.Contents: %w", err)
	}

	return c.mode.Apply(keys, keysValue, c.extractor)
}
