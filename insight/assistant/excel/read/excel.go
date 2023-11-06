package read

import (
	"errors"
	"fmt"

	"github.com/xuri/excelize/v2"
)

type Excel struct {
	path      string
	excelFile *excelize.File
}

func NewExcel(path string) (*Excel, error) {
	e := &Excel{}
	e.path = path

	var err error
	e.excelFile, err = excelize.OpenFile(e.path)
	if err != nil {
		return nil, fmt.Errorf("NewExcel error; %w", err)
	}

	return e, nil
}

func (e *Excel) readSheet(config Config) ([][]string, error) {
	if config.SheetName == "" {
		if config.SheetIndex <= 0 {
			return nil, errors.New("sheet name or index no exists")
		}

		sheetList := e.excelFile.GetSheetList()
		config.SheetName = sheetList[config.SheetIndex-1]
	}

	if config.SheetName == "" {
		return nil, errors.New("sheet name or index no exists")
	}

	rows, err := e.excelFile.GetRows(config.SheetName)
	if err != nil {
		return nil, fmt.Errorf("GetRows error; %w", err)
	}

	if config.StartRow > 1 {
		rows = rows[config.StartRow-1:]
	}

	if len(config.ColsIndex) > 0 {
		var newRows [][]string
		for _, row := range rows {
			var newRow []string
			for _, index := range config.ColsIndex {
				newRow = append(newRow, row[index])
			}

			newRows = append(newRows, newRow)
		}

		rows = newRows
	}

	return rows, nil
}

func (e *Excel) Close() error {
	if err := e.excelFile.Close(); err != nil {
		return fmt.Errorf("close excel error; %w", err)
	}

	return nil
}
