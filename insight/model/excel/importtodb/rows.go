package importtodb

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/model/excel/read"
	"github.com/auho/go-etl/v2/insight/model/table"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ sourceor = (*RowsImportToDb)(nil)

type RowsImportToDb struct {
	Source
	Titles    []string // save to db 的 columns
	RowsTable *table.RowsTable
}

func (ri *RowsImportToDb) GetTable() table.Tabler {
	return ri.RowsTable
}

func (ri *RowsImportToDb) GetTitles() []string {
	return ri.Titles
}

func (ri *RowsImportToDb) GetSheetData() (read.SheetDataor, error) {
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

func (ri *RowsImportToDb) Import(db *simpleDb.SimpleDB) error {
	return RunImportToDb(db, ri)
}
