package load

import (
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/excel/reader"
	"github.com/auho/go-etl/v3/insight/assistant/schema/create"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ Resource = (*Raw)(nil)

type Raw struct {
	baseResource

	Rows      assistant.Entity
	sheetData *reader.SheetDataWithTitle
}

func (r *Raw) DB() *simpledb.SimpleDB {
	return r.Rows.DB()
}

func (r *Raw) Prepare() error {
	return nil
}

func (r *Raw) Name() string {
	return r.Rows.Name()
}

func (r *Raw) Tabler() create.Tabler {
	return create.NewRowsTable(r.Rows)
}

func (r *Raw) TitlesName() []string {
	return r.sheetData.GetTitles()
}

func (r *Raw) TitlesIndex() []int {
	var indexes []int
	for i := range r.TitlesName() {
		indexes = append(indexes, i)
	}

	return indexes
}

func (r *Raw) SheetData(excel *reader.Excel) (reader.SheetDataReader, error) {
	var err error
	r.sheetData, err = r.readSheetData(excel, r.buildSheetConfig())

	return r.sheetData, err
}

func (r *Raw) readSheetData(excel *reader.Excel, sheetConfig reader.Config) (*reader.SheetDataWithTitle, error) {
	sheetData, err := reader.NewSheetDataWithTitle(excel, sheetConfig, nil)
	if err != nil {
		return nil, fmt.Errorf("NewSheetDataWithTitle: %w", err)
	}

	err = sheetData.ReadData()
	if err != nil {
		return nil, fmt.Errorf("ReadData: %w", err)
	}

	return sheetData, nil
}
