package table

import (
	"github.com/auho/go-etl/v2/insight/model"
)

type DataTable struct {
	table
	data model.Dataor
}

func NewDataTable(data model.Dataor) *DataTable {
	t := &DataTable{}
	t.data = data

	t.buildData()

	return t
}

func (t *DataTable) buildData() {
	t.initCommand(t.data.TableName())
}

func (t *DataTable) BuildDataForTag(command *command) {
	command.AddKeyBigInt(t.data.GetIdName())
}
