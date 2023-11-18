package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
)

type RawTable struct {
	table
	raw assistant.Rawer
}

func NewRawTable(raw assistant.Rawer, opts ...TableOption) *RawTable {
	t := &RawTable{}
	t.raw = raw
	t.db = raw.GetDB()

	t.options(opts)
	t.build()

	return t
}

func (t *RawTable) build() {
	t.initCommand(t.raw.TableName())

	t.execRawCommandFunc(t.raw)
}

func (t *RawTable) WithCommand(fn func(*tablestructure.Command)) *RawTable {
	fn(t.Command)

	return t
}
