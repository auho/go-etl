package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/schema"
)

type DataTable struct {
	table
	data assistant.Entity
}

func NewDataTable(data assistant.Entity, opts ...TableOption) *DataTable {
	t := &DataTable{}
	t.data = data
	t.db = t.data.GetDB()

	t.options(opts)
	t.build()

	return t
}

func (t *DataTable) build() {
	t.initCommand(t.data.TableName())
	t.AddPKBigInt(t.data.GetIDName())

	t.execRawCommandFunc(t.data)
}

func (t *DataTable) BuildForTag(command *schema.Command) {
	command.AddKeyBigInt(t.data.GetIDName())
}

func (t *DataTable) WithCommand(fn func(*schema.Command)) *DataTable {
	fn(t.Command)

	return t
}
