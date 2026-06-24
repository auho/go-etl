package dbimport

import (
	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/excel/read"
	"github.com/auho/go-etl/v3/insight/assistant/schema/buildtable"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ Resource = (*RowsResource)(nil)

type RowsResource struct {
	ResourceBase
	Titles // column title of save to db
	Rows   assistant.Entity
}

func (rs *RowsResource) Prepare() error {
	return rs.Titles.prepare()
}

func (rs *RowsResource) GetName() string {
	return rs.Rows.GetName()
}

func (rs *RowsResource) GetTable() buildtable.Tabler {
	return buildtable.NewRowsTable(rs.Rows)
}

func (rs *RowsResource) GetSheetData(excel *read.Excel) (read.SheetDataReader, error) {
	return rs.readSheetData(excel, rs.buildSheetConfig())
}

func (rs *RowsResource) GetDB() *simpledb.SimpleDB {
	return rs.Rows.GetDB()
}
