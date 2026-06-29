package collector

import (
	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collector/mode"
	"github.com/auho/go-etl/v3/job/transform/collector/source/keys"
)

// This file provides convenience constructors that combine a Keys source with
// a mode, so callers can build a Collector in one call without importing the
// keys and mode packages separately.
//
// Two selection semantics are offered:
//   - Semantic A (by key position): selects keys by their position in the keys
//     list, regardless of whether the value is empty.
//   - Semantic B (by value, "Valued" suffix): selects only keys with non-empty
//     values, then applies positional logic.

// --- Semantic A: by key position ---

// NewKeysAll searches the contents of ALL keys together.
func NewKeysAll(ks []string, e extract.Extractor) *Collector {
	return NewCollector(keys.New(ks), mode.NewAll(), e)
}

// NewKeysFirst searches only the content of the first key (keys[0]).
func NewKeysFirst(ks []string, e extract.Extractor) *Collector {
	return NewCollector(keys.New(ks), mode.NewFirst(), e)
}

// NewKeysLast searches only the content of the last key.
func NewKeysLast(ks []string, e extract.Extractor) *Collector {
	return NewCollector(keys.New(ks), mode.NewLast(), e)
}

// NewKeysFirstN searches the contents of the first n keys.
func NewKeysFirstN(ks []string, n int, e extract.Extractor) *Collector {
	return NewCollector(keys.New(ks), mode.NewFirstN(n), e)
}

// NewKeysLastN searches the contents of the last n keys.
func NewKeysLastN(ks []string, n int, e extract.Extractor) *Collector {
	return NewCollector(keys.New(ks), mode.NewLastN(n), e)
}

// NewKeysMatchAny searches each key's content one by one, stopping at the
// first match. Empty-value keys are still searched.
func NewKeysMatchAny(ks []string, e extract.Extractor) *Collector {
	return NewCollector(keys.New(ks), mode.NewMatchAny(), e)
}

// NewKeysMatchAnyN searches the first n keys one by one, stopping at the
// first match. Empty-value keys are still searched.
func NewKeysMatchAnyN(ks []string, n int, e extract.Extractor) *Collector {
	return NewCollector(keys.New(ks), mode.NewMatchAnyN(n), e)
}

// --- Semantic B: by value (only non-empty valued keys) ---

// NewKeysAllValued searches the contents of all keys with non-empty values.
// Empty-value keys are skipped.
func NewKeysAllValued(ks []string, e extract.Extractor) *Collector {
	return NewCollector(keys.New(ks), mode.NewAllValued(), e)
}

// NewKeysFirstValued searches the content of the first key with a non-empty
// value.
func NewKeysFirstValued(ks []string, e extract.Extractor) *Collector {
	return NewCollector(keys.New(ks), mode.NewFirstValued(), e)
}

// NewKeysLastValued searches the content of the last key with a non-empty
// value.
func NewKeysLastValued(ks []string, e extract.Extractor) *Collector {
	return NewCollector(keys.New(ks), mode.NewLastValued(), e)
}

// NewKeysFirstValuedN searches the contents of the first n keys with
// non-empty values.
func NewKeysFirstValuedN(ks []string, n int, e extract.Extractor) *Collector {
	return NewCollector(keys.New(ks), mode.NewFirstValuedN(n), e)
}

// NewKeysLastValuedN searches the contents of the last n keys with non-empty
// values.
func NewKeysLastValuedN(ks []string, n int, e extract.Extractor) *Collector {
	return NewCollector(keys.New(ks), mode.NewLastValuedN(n), e)
}

// NewKeysMatchAnyValued iterates only keys with non-empty values, searching
// each one until a match is found. Empty-value keys are skipped entirely.
func NewKeysMatchAnyValued(ks []string, e extract.Extractor) *Collector {
	return NewCollector(keys.New(ks), mode.NewMatchAnyValued(), e)
}

// NewKeysMatchAnyValuedN iterates the first n keys with non-empty values,
// searching each one until a match is found.
func NewKeysMatchAnyValuedN(ks []string, n int, e extract.Extractor) *Collector {
	return NewCollector(keys.New(ks), mode.NewMatchAnyValuedN(n), e)
}
