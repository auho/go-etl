package read

import (
	"github.com/auho/go-etl/v2/insight/assistant/model"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure/buildtable"
)

var _execl *Excel
var _raw *model.Raw

func ExampleNewSchema() {
	s, _ := NewSchema(_execl, buildtable.NewRawTable(_raw), Config{
		SheetName:  "",
		SheetIndex: 0,
		StartRow:   0,
		EndRow:     0,
		ColsIndex:  nil,
	})

	// handler title func
	_ = s.WithFuncTitle(func(title string) string {
		// handler title

		return title
	})

	// build table => build table.rawTable
	_, _ = s.BuildTable()
}

func ExampleNewSchemaWithPath() {
	s, _ := NewSchemaWithPath("xlsxPath", buildtable.NewRawTable(_raw), Config{
		SheetName:  "",
		SheetIndex: 0,
		StartRow:   0,
		EndRow:     0,
		ColsIndex:  nil,
	})

	// handler title func
	_ = s.WithFuncTitle(func(title string) string {
		// handler title

		return title
	})

	// build table => build table.rawTable
	_, _ = s.BuildTable()
}

func ExampleNewSheetDataNoTitle() {
	s, _ := NewSheetDataNoTitle(
		_execl,
		Config{
			SheetName:  "",
			SheetIndex: 0,
			StartRow:   0,
			EndRow:     0,
			ColsIndex:  nil,
		})

	// read data
	_ = s.ReadData()

	// handler rows
	_ = s.HandlerRows(func(rows [][]string) ([][]string, error) {
		return rows, nil
	})

	// get rows
	_ = s.GetRows()

	// get rows with any
	_ = s.GetRowsWithAny()
}

func ExampleNewSheetDataNoTitleWithPath() {
	s, _ := NewSheetDataNoTitleWithPath(
		"xlsxPath",
		Config{
			SheetName:  "",
			SheetIndex: 0,
			StartRow:   0,
			EndRow:     0,
			ColsIndex:  nil,
		})

	// read data
	_ = s.ReadData()

	// handler rows
	_ = s.HandlerRows(func(rows [][]string) ([][]string, error) {
		return rows, nil
	})

	// get rows
	_ = s.GetRows()

	// get rows with any
	_ = s.GetRowsWithAny()
}

func ExampleNewSheetDataWithTitle() {
	s, _ := NewSheetDataWithTitle(
		_execl,
		Config{
			SheetName:  "",
			SheetIndex: 0,
			StartRow:   0,
			EndRow:     0,
			ColsIndex:  nil,
		},
		map[string]string{"title1": "title1_alias"},
	)

	// read data
	_ = s.ReadData()

	// handler rows
	_ = s.HandlerRows(func(rows [][]string) ([][]string, error) {
		return rows, nil
	})

	// get titles
	_ = s.GetTitles()

	// get alias
	_ = s.GetAlias()

	// get rows
	_ = s.GetRows()

	// get rows with any
	_ = s.GetRowsWithAny()
}

func ExampleNewSheetDataWithTitleWithPath() {
	s, _ := NewSheetDataWithTitleWithPath(
		"xlsxPath",
		Config{
			SheetName:  "",
			SheetIndex: 0,
			StartRow:   0,
			EndRow:     0,
			ColsIndex:  nil,
		},
		map[string]string{"title1": "title1_alias"},
	)

	// read data
	_ = s.ReadData()

	// handler rows
	_ = s.HandlerRows(func(rows [][]string) ([][]string, error) {
		return rows, nil
	})

	// get titles
	_ = s.GetTitles()

	// get alias
	_ = s.GetAlias()

	// get rows
	_ = s.GetRows()

	// get rows with any
	_ = s.GetRowsWithAny()
}
