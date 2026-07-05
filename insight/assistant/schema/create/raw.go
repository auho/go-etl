package create

import (
	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
)

type RawTable struct {
	table
	raw assistant.Raw
}

func NewRawTable(raw assistant.Raw, opts ...TableOption) *RawTable {
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

func (t *RawTable) WithCommand(fn func(*schema.Command)) *RawTable {
	fn(t.Command)

	return t
}
