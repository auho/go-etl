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

	rowsScan, err := e.excelFile.Rows(config.SheetName)
	if err != nil {
		return nil, fmt.Errorf("rows error; %w", err)
	}

	var _i = 0
	var rows [][]string
	for rowsScan.Next() {
		_i += 1
		if _i < config.StartRow {
			continue
		}

		// end row > 0 AND current > end row
		if config.EndRow > 0 && _i > config.EndRow {
			break
		}

		row, err1 := rowsScan.Columns()
		if err1 != nil {
			return nil, fmt.Errorf("rows scan columns error; %w", err)
		}

		rows = append(rows, row)
	}

	if err = rowsScan.Close(); err != nil {
		return nil, fmt.Errorf("rows scan close error; %w", err)
	}

	if len(config.ColsIndex) > 0 {
		var newRows [][]string
		for _, row := range rows {
			if len(row) <= 0 {
				continue
			}

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
