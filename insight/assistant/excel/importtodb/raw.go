package importtodb

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/excel/read"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure/buildtable"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ Resourcer = (*RawResource)(nil)

type RawResource struct {
	Resource

	Rows assistant.Rowsor

	titlesName []string
	sheetData  *read.SheetDataWithTitle
}

func (rs *RawResource) GetDB() *simpleDb.SimpleDB {
	return rs.Rows.GetDB()
}

func (rs *RawResource) Prepare() error {
	return nil
}

func (rs *RawResource) GetName() string {
	return rs.Rows.GetName()
}

func (rs *RawResource) GetTable() buildtable.Tabler {
	return buildtable.NewRowsTable(rs.Rows)
}

func (rs *RawResource) GetTitlesName() []string {
	return rs.sheetData.GetTitles()
}

func (rs *RawResource) GetTitlesIndex() []int {
	var indexes []int
	for i := range rs.GetTitlesName() {
		indexes = append(indexes, i)
	}

	return indexes
}

func (rs *RawResource) GetSheetData(excel *read.Excel) (read.SheetDataor, error) {
	var err error
	rs.sheetData, err = rs.readSheetData(excel, rs.buildSheetConfig())

	return rs.sheetData, err
}

func (rs *RawResource) readSheetData(excel *read.Excel, sheetConfig read.Config) (*read.SheetDataWithTitle, error) {
	sheetData, err := read.NewSheetDataWithTitle(excel, sheetConfig, nil)
	if err != nil {
		return nil, fmt.Errorf("NewSheetDataWithTitle error; %w", err)
	}

	err = sheetData.ReadData()
	if err != nil {
		return nil, fmt.Errorf("ReadData error; %w", err)
	}

	return sheetData, nil
}
