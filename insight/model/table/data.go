package table

import (
	"github.com/auho/go-etl/v2/insight/model"
)

type DataTable struct {
	table
	data *model.Data
}

func NewDataTable(data *model.Data) *DataTable {
	dt := &DataTable{}
	dt.data = data

	dt.buildData()

	return dt
}

func (dt *DataTable) buildData() {
	dt.initTable(dt.data.TableName())
}
