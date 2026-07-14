package reader

import (
	"github.com/auho/go-etl/v3/insight/assistant/entity"
	"github.com/auho/go-etl/v3/insight/assistant/schema/create"
)

var _excel *Excel
var _raw *entity.Raw

func ExampleNewSchema() {
	s, _ := NewSchema(_excel, create.NewRawTable(_raw), Config{
		SheetName:     "",
		SheetIndex:    0,
		StartRow:      0,
		EndRow:        0,
		ColumnIndexes: nil,
	})

	// handler title func
	_ = s.WithTitleFunc(func(title string) string {
		// handler title

		return title
	})

	// build table => build table.rawTable
	_, _ = s.BuildTable()
}

func ExampleNewSchemaWithPath() {
	s, _ := NewSchemaWithPath("xlsxPath", create.NewRawTable(_raw), Config{
		SheetName:     "",
		SheetIndex:    0,
		StartRow:      0,
		EndRow:        0,
		ColumnIndexes: nil,
	})

	// handler title func
	_ = s.WithTitleFunc(func(title string) string {
		// handler title

		return title
	})

	// build table => build table.rawTable
	_, _ = s.BuildTable()
}

func ExampleNewSheetDataNoTitle() {
	s, _ := NewSheetDataNoTitle(
		_excel,
		Config{
			SheetName:     "",
			SheetIndex:    0,
			StartRow:      0,
			EndRow:        0,
			ColumnIndexes: nil,
		})

	// read data
	_ = s.ReadData()

	// transform rows
	_ = s.TransformRows(func(rows [][]string) ([][]string, error) {
		return rows, nil
	})

	// get rows
	_ = s.GetRows()

	// get rows as any
	_ = s.GetRowsAsAny()
}

func ExampleNewSheetDataNoTitleWithPath() {
	s, _ := NewSheetDataNoTitleWithPath(
		"xlsxPath",
		Config{
			SheetName:     "",
			SheetIndex:    0,
			StartRow:      0,
			EndRow:        0,
			ColumnIndexes: nil,
		})

	// read data
	_ = s.ReadData()

	// transform rows
	_ = s.TransformRows(func(rows [][]string) ([][]string, error) {
		return rows, nil
	})

	// get rows
	_ = s.GetRows()

	// get rows as any
	_ = s.GetRowsAsAny()
}

func ExampleNewSheetDataWithTitle() {
	s, _ := NewSheetDataWithTitle(
		_excel,
		Config{
			SheetName:     "",
			SheetIndex:    0,
			StartRow:      0,
			EndRow:        0,
			ColumnIndexes: nil,
		},
		map[string]string{"title1": "title1_alias"},
	)

	// read data
	_ = s.ReadData()

	// transform rows
	_ = s.TransformRows(func(rows [][]string) ([][]string, error) {
		return rows, nil
	})

	// get titles
	_ = s.GetTitles()

	// get alias
	_ = s.GetAlias()

	// get rows
	_ = s.GetRows()

	// get rows as any
	_ = s.GetRowsAsAny()
}

func ExampleNewSheetDataWithTitleWithPath() {
	s, _ := NewSheetDataWithTitleWithPath(
		"xlsxPath",
		Config{
			SheetName:     "",
			SheetIndex:    0,
			StartRow:      0,
			EndRow:        0,
			ColumnIndexes: nil,
		},
		map[string]string{"title1": "title1_alias"},
	)

	// read data
	_ = s.ReadData()

	// transform rows
	_ = s.TransformRows(func(rows [][]string) ([][]string, error) {
		return rows, nil
	})

	// get titles
	_ = s.GetTitles()

	// get alias
	_ = s.GetAlias()

	// get rows
	_ = s.GetRows()

	// get rows as any
	_ = s.GetRowsAsAny()
}
