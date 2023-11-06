package importtodb

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant/excel/read"
	buildtable2 "github.com/auho/go-etl/v2/insight/assistant/tablestructure/buildtable"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ resourcer = (*RowsResource)(nil)

type RowsResource struct {
	Resource
	Titles    []string // save to db 的 columns
	RowsTable *buildtable2.RowsTable
}

func (ri *RowsResource) GetTable() buildtable2.Tabler {
	return ri.RowsTable
}

func (ri *RowsResource) GetTitles() []string {
	return ri.Titles
}

func (ri *RowsResource) GetSheetData() (read.SheetDataor, error) {
	sheetData, err := read.NewSheetDataNoTitle(ri.XlsxPath, ri.SheetName, ri.StartRow)
	if err != nil {
		return nil, fmt.Errorf("NewSheetDataNoTitle error; %w", err)
	}

	err = sheetData.ReadData()
	if err != nil {
		return nil, fmt.Errorf("ReadData error; %w", err)
	}

	return sheetData, nil
}

func (ri *RowsResource) Import(db *simpleDb.SimpleDB) error {
	return RunImportToDb(db, ri)
}
