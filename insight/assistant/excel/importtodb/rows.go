package importtodb

import (
	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/excel/read"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure/buildtable"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ Resourcer = (*RowsResource)(nil)

type RowsResource struct {
	Resource
	Titles
	Rows       assistant.Rowsor
	titlesKey  []string
	titleIndex []int
}

func (rs *RowsResource) Prepare() error {
	return rs.Titles.prepare()
}

func (rs *RowsResource) GetTable() buildtable.Tabler {
	return buildtable.NewRowsTable(rs.Rows)
}

func (rs *RowsResource) GetSheetData(excel *read.Excel) (read.SheetDataor, error) {
	return rs.readSheetData(excel, rs.buildSheetConfig())
}

func (rs *RowsResource) GetDB() *simpleDb.SimpleDB {
	return rs.Rows.GetDB()
}
