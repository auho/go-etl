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

	t.buildData()

	return t
}

func (t *DataTable) buildData() {
	t.initCommand(t.data.TableName())
	t.AddPkBigInt(t.data.GetIdName())
}

func (t *DataTable) BuildDataForTag(command *tablestructure.Command) {
	command.AddKeyBigInt(t.data.GetIdName())
}

func (t *DataTable) Exec(fn func(command *tablestructure.Command)) *DataTable {
	fn(t.Command)

	return t
}
