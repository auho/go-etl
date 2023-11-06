package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
)

type DataTable struct {
	table
	data assistant.Dataor
}

func NewDataTable(data assistant.Dataor) *DataTable {
	t := &DataTable{}
	t.data = data

	t.build()

	return t
}

func (t *DataTable) build() {
	t.initCommand(t.data.TableName())
	t.AddPkBigInt(t.data.GetIdName())
}

func (t *DataTable) BuildForTag(command *tablestructure.Command) {
	command.AddKeyBigInt(t.data.GetIdName())
}

func (t *DataTable) ExecCommand(fn func(command *tablestructure.Command)) *DataTable {
	fn(t.Command)

	return t
}
