package read

import (
	"fmt"
)

type SheetDataor interface {
	ReadData() error
	HandlerRows(fn func(rows [][]string) ([][]string, error)) error
	GetRows() [][]string
	GetRowsWithAny() [][]any
}

type SheetData struct {
	xlsxPath  string
	sheetName string
	startRow  int
	rows      [][]string
}

func (sd *SheetData) readFromSheet() error {
	excel, err := NewExcel(sd.xlsxPath)
	if err != nil {
		return fmt.Errorf("NewExcel error; %w", err)
	}

	sd.rows, err = excel.excelFile.GetRows(sd.sheetName)
	if err != nil {
		return fmt.Errorf("GetRows error; %w", err)
	}

	if sd.startRow > 0 {
		sd.rows = sd.rows[sd.startRow:]
	}

	return nil
}

func (sd *SheetData) GetRows() [][]string {
	return sd.rows
}

func (sd *SheetData) GetRowsWithAny() [][]any {
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

func (sd *SheetData) HandlerRows(fn func(rows [][]string) ([][]string, error)) error {
	var err error
	sd.rows, err = fn(sd.rows)
	if err != nil {
		return fmt.Errorf("HandlerRows fn errro; %w", err)
	}

	return nil
}
