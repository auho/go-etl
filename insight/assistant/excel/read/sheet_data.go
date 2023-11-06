package read

import (
	"fmt"
)

type sheetData struct {
	excel  *Excel
	config Config
	rows   [][]string
}

func (sd *sheetData) readSheet() error {
	var err error
	sd.rows, err = sd.excel.readSheet(sd.config)
	return err
}

func (sd *sheetData) GetRows() [][]string {
	return sd.rows
}

func (sd *sheetData) GetRowsWithAny() [][]any {
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

func (sd *sheetData) HandlerRows(fn func(rows [][]string) ([][]string, error)) error {
	var err error
	sd.rows, err = fn(sd.rows)
	if err != nil {
		return fmt.Errorf("HandlerRows fn errro; %w", err)
	}

	return nil
}
