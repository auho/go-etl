package dataset

import (
	"fmt"
)

type Mode string

const ModeAppend Mode = "append"
const ModeSpread Mode = "spread"

// Moder
// how to merge datasets
type Moder interface {
	Data() (*Data, error)
	Name() string
	Sets() []Set
}

func NewMode(mode Mode, ds *Dataset) (Moder, error) {
	var dsMode Moder
	switch mode {
	case ModeAppend:
		dsMode = NewAppendMode(ds)
	case ModeSpread:
		dsMode = NewSpreadMode(ds)
	default:
		return nil, fmt.Errorf("dataset mode[%s] error", mode)
	}

	return dsMode, nil
}
