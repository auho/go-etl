package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant"
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

func (r *RawTable) build() {
	r.initCommand(r.raw.TableName())

	r.execCommandFunc()
	r.execRawCommandFunc(r.raw)
}
