package collector

import (
	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-etl/v3/job/transform/collector/keys"
	"github.com/auho/go-etl/v3/job/transform/collector/mode"
)

// NewKeysAll creates a Collector with a Keys source and All mode.
func NewKeysAll(ks []string, e extract.Extractor) *Collector {
	return NewCollector(keys.New(ks), mode.NewAll(), e)
}

// NewKeysFirst creates a Collector with a Keys source and First mode.
func NewKeysFirst(ks []string, e extract.Extractor) *Collector {
	return NewCollector(keys.New(ks), mode.NewFirst(), e)
}

// NewKeysLast creates a Collector with a Keys source and Last mode.
func NewKeysLast(ks []string, e extract.Extractor) *Collector {
	return NewCollector(keys.New(ks), mode.NewLast(), e)
}

// NewKeysFirstN creates a Collector with a Keys source and FirstN mode.
func NewKeysFirstN(ks []string, n int, e extract.Extractor) *Collector {
	return NewCollector(keys.New(ks), mode.NewFirstN(n), e)
}

// NewKeysLastN creates a Collector with a Keys source and LastN mode.
func NewKeysLastN(ks []string, n int, e extract.Extractor) *Collector {
	return NewCollector(keys.New(ks), mode.NewLastN(n), e)
}

// NewKeysMatchAny creates a Collector with a Keys source and MatchAny mode.
func NewKeysMatchAny(ks []string, e extract.Extractor) *Collector {
	return NewCollector(keys.New(ks), mode.NewMatchAny(), e)
}

// NewKeysMatchAnyN creates a Collector with a Keys source and MatchAnyN mode.
func NewKeysMatchAnyN(ks []string, n int, e extract.Extractor) *Collector {
	return NewCollector(keys.New(ks), mode.NewMatchAnyN(n), e)
}
