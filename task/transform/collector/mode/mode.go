package mode

import (
	"github.com/auho/go-etl/v3/task/extract"
)

// Mode decides which keys to select and how to drive the Extractor.
// It receives the ordered keys and a key-to-content map instead of a Source,
// so mode has no dependency on the source package.
//
// There are two selection semantics:
//   - Semantic A (by key position): selects keys by their position in the keys
//     list, regardless of whether the value is empty. e.g. First, Last, FirstN.
//   - Semantic B (by value): selects only keys with non-empty values, then
//     applies positional logic. e.g. FirstValued, LastValued, FirstValuedN.
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

// valuedKeys returns keys whose value is non-empty, preserving order.
func valuedKeys(keys []string, keysValue map[string]string) []string {
	valued := make([]string, 0, len(keys))
	for _, k := range keys {
		if keysValue[k] != "" {
			valued = append(valued, k)
		}
	}
	return valued
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
