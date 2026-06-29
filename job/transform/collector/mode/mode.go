package mode

import (
	"fmt"

	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collector"
)

// searchContents fetches contents for the given keys and runs a single Search.
func searchContents(source collector.Source, ks []string, item map[string]any, e extract.Extractor) (extract.Result, error) {
	contents, err := source.Contents(ks, item)
	if err != nil {
		return extract.Result{}, fmt.Errorf("contents: %w", err)
	}
	return e.Search(contents), nil
}

// takeFirst returns the first n elements of keys; if n >= len, returns all.
func takeFirst(keys []string, n int) []string {
	if n >= len(keys) {
		return keys
	}
	return keys[:n]
}

// takeLast returns the last n elements of keys; if n >= len, returns all.
func takeLast(keys []string, n int) []string {
	if n >= len(keys) {
		return keys
	}
	return keys[len(keys)-n:]
}