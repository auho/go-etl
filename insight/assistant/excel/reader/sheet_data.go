package reader

import (
	"fmt"
)

// sheetData holds the shared state and behavior embedded by SheetDataNoTitle and SheetDataWithTitle.
type sheetData struct {
	excel  *Excel
	config Config
	rows   [][]string
}

// readSheet reads rows from the Excel file using the configured Config.
func (sd *sheetData) readSheet() error {
	var err error
	sd.rows, err = sd.excel.readSheet(sd.config)
	return err
}

// GetRows returns the current rows as [][]string.
func (sd *sheetData) GetRows() [][]string {
	return sd.rows
}

// GetRowsAsAny returns the current rows as [][]any for use with generic insert APIs.
func (sd *sheetData) GetRowsAsAny() [][]any {
	var data [][]any
	for _, row := range sd.rows {
		var rowAny []any
		for _, value := range row {
			rowAny = append(rowAny, value)
		}

		data = append(data, rowAny)
	}

	return data
}

// TransformRows applies fn to the current rows, replacing them with the result.
func (sd *sheetData) TransformRows(fn func(rows [][]string) ([][]string, error)) error {
	var err error
	sd.rows, err = fn(sd.rows)
	if err != nil {
		return fmt.Errorf("fn: %w", err)
	}

	return nil
}
