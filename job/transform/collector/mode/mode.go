package mode

import (
	"github.com/auho/go-etl/v3/job/extract"
)

// Mode decides which keys to select and how to drive the Extractor.
// It receives the ordered keys and a key-to-content map instead of a Source,
// so mode has no dependency on the source package.
type Mode interface {
	Prepare() error
	Apply(keys []string, keysValue map[string]string, e extract.Extractor) (extract.Result, error)
}

// valuesByKeys collects contents for the given keys in order.
func valuesByKeys(keys []string, keysValue map[string]string) []string {
	values := make([]string, len(keys))
	for i, k := range keys {
		values[i] = keysValue[k]
	}
	return values
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
