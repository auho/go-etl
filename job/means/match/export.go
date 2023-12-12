package match

import (
	"github.com/auho/go-etl/v2/job/means"
)

type Export struct {
	Ok   bool
	Rule means.Ruler

	Keys          []string
	DefaultValues map[string]any
}

func (e *Export) IsOk() bool {
	return e.Ok
}

func (e *Export) GetKeys() []string {
	return e.Keys
}

func (e *Export) GetDefaultValues() map[string]any {
	return e.DefaultValues
}
