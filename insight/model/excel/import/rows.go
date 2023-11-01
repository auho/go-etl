package _import

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/model/excel/read"
	"github.com/auho/go-etl/v2/insight/model/table"
)

var _ sourceor = (*RowsImport)(nil)

type RowsImport struct {
	Source
	Titles    []string // save to db 的 columns
	RowsTable *table.RowsTable
}

func (ri *RowsImport) GetTable() table.Tabler {
	return ri.RowsTable
}

func (ri *RowsImport) GetTitles() []string {
	return ri.Titles
}

func (ri *RowsImport) GetSheetData() (read.SheetDataor, error) {
	sheetData, err := read.NewSheetDataNoTitle(ri.XlsxPath, ri.SheetName, ri.StartRow)
	if err != nil {
		return nil, fmt.Errorf("NewSheetDataNoTitle error; %w", sheetData)
	}

	err = sheetData.ReadData()
	if err != nil {
		return nil, fmt.Errorf("ReadData error; %w", err)
	}

	return sheetData, nil
}
