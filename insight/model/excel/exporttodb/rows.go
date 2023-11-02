package exporttodb

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/model/excel/read"
	"github.com/auho/go-etl/v2/insight/model/table"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ sourceor = (*RowsExportToDb)(nil)

type RowsExportToDb struct {
	Source
	Titles    []string // save to db 的 columns
	RowsTable *table.RowsTable
}

func (re *RowsExportToDb) GetTable() table.Tabler {
	return re.RowsTable
}

func (re *RowsExportToDb) GetTitles() []string {
	return re.Titles
}

func (re *RowsExportToDb) GetSheetData() (read.SheetDataor, error) {
	sheetData, err := read.NewSheetDataNoTitle(re.XlsxPath, re.SheetName, re.StartRow)
	if err != nil {
		return nil, fmt.Errorf("NewSheetDataNoTitle error; %w", err)
	}

	err = sheetData.ReadData()
	if err != nil {
		return nil, fmt.Errorf("ReadData error; %w", err)
	}

	return sheetData, nil
}

func (re *RowsExportToDb) Export(db *simpleDb.SimpleDB) error {
	return RunExportToDb(db, re)
}
