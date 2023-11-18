package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
)

type DataTable struct {
	table
	data assistant.Dataor
}

func NewDataTable(data assistant.Dataor, opts ...TableOption) *DataTable {
	t := &DataTable{}
	t.data = data
	t.db = t.data.GetDB()

	t.options(opts)
	t.build()

	return t
}

func (t *DataTable) build() {
	t.initCommand(t.data.TableName())
	t.AddPkBigInt(t.data.GetIdName())

	t.execRawCommandFunc(t.data)
}

func (t *DataTable) BuildForTag(command *tablestructure.Command) {
	command.AddKeyBigInt(t.data.GetIdName())
}

func (t *DataTable) WithCommand(fn func(*tablestructure.Command)) *DataTable {
	fn(t.Command)

	return t
}
