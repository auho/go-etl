package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
)

type TableOption func(Tabler)

func WithCommand(fn func(command *tablestructure.Command)) func(Tabler) {
	return func(t Tabler) {
		t.withCommand(fn)
	}
}

func WithConfig(config Config) func(Tabler) {
	return func(t Tabler) {
		t.withConfig(config)
	}
}
