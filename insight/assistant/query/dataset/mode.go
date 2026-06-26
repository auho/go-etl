package dataset

import (
	"fmt"
)

type MergeMode string

const ModeAppend MergeMode = "append"
const ModeSpread MergeMode = "spread"

// Mode
// how to merge datasets
type Mode interface {
	Data() (*Data, error)
	Name() string
	Sets() []Set
}

func NewMode(mode MergeMode, ds *Dataset) (Mode, error) {
	var dsMode Mode
	switch mode {
	case ModeAppend:
		dsMode = NewAppendMode(ds)
	case ModeSpread:
		dsMode = NewSpreadMode(ds)
	default:
		return nil, fmt.Errorf("mode[%s] is invalid", mode)
	}

	return dsMode, nil
}
