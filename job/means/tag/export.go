package tag

import (
	"github.com/auho/go-etl/v2/job/means"
)

type export struct {
	Ok   bool
	Rule means.Ruler

	keys          []string
	defaultValues map[string]any
}

func (e *export) IsOk() bool {
	return e.Ok
}

func (e *export) GetKeys() []string {
	return e.keys
}

func (e *export) DefaultValues() map[string]any {
	return e.defaultValues
}
