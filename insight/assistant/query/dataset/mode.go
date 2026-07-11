package dataset

import (
	"fmt"
)

// MergeMode defines how datasets are merged.
type MergeMode string

const MergeModeAppend MergeMode = "append"
const MergeModeSpread MergeMode = "spread"

// Merger defines how to merge datasets.
type Merger interface {
	Data() (*Result, error)
	Name() string
	Sets() []Subset
}

// NewMerger creates a Merger based on the given merge mode.
func NewMerger(mode MergeMode, ds *Dataset) (Merger, error) {
	if ds == nil {
		return nil, fmt.Errorf("dataset is nil")
	}

	var m Merger
	switch mode {
	case MergeModeAppend:
		m = NewAppendMode(ds)
	case MergeModeSpread:
		m = NewSpreadMode(ds)
	default:
		return nil, fmt.Errorf("invalid merge mode: %s", mode)
	}

	return m, nil
}
