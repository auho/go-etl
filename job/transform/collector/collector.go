package collector

import (
	"fmt"

	"github.com/auho/go-etl/v3/job/extract"
)

// Source provides ordered content values from a data source.
type Source interface {
	Title() string
	Keys() []string
	Prepare() error
	Contents(keys []string, item map[string]any) ([]string, error)
}

// Mode decides which keys to select and how to drive the Extractor.
type Mode interface {
	Prepare() error
	Apply(source Source, item map[string]any, e extract.Extractor) (extract.Result, error)
}

// Collector orchestrates a Source, a Mode, and an Extractor.
// All lifecycle and extraction logic is unified here.
type Collector struct {
	source    Source
	mode      Mode
	extractor extract.Extractor
}

func NewCollector(source Source, mode Mode, extractor extract.Extractor) *Collector {
	return &Collector{source: source, mode: mode, extractor: extractor}
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
		return err
	}
	if err := c.mode.Prepare(); err != nil {
		return err
	}
	return c.extractor.Prepare()
}

func (c *Collector) Close() error {
	return c.extractor.Close()
}

func (c *Collector) Extract(item map[string]any) (extract.Result, error) {
	return c.mode.Apply(c.source, item, c.extractor)
}