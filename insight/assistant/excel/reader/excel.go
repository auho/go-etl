package reader

import (
	"errors"
	"fmt"

	"github.com/xuri/excelize/v2"
)

// Excel wraps an excelize.File for reading sheet data.
type Excel struct {
	path      string
	excelFile *excelize.File
}

// NewExcel opens an xlsx file at path and returns an Excel ready for reading.
func NewExcel(path string) (*Excel, error) {
	e := &Excel{}
	e.path = path

	var err error
	e.excelFile, err = excelize.OpenFile(e.path)
	if err != nil {
		return nil, fmt.Errorf("OpenFile: %w", err)
	}

	return e, nil
}

// readSheet reads rows from the sheet identified by config.
// If ColumnIndexes is non-empty, only those columns are kept; missing columns default to "".
func (e *Excel) readSheet(config Config) ([][]string, error) {
	if config.SheetName == "" {
		if config.SheetIndex <= 0 {
			return nil, errors.New("sheet name or index does not exist")
		}

		sheetList := e.excelFile.GetSheetList()
		config.SheetName = sheetList[config.SheetIndex-1]
	}

	if config.SheetName == "" {
		return nil, errors.New("sheet name or index does not exist")
	}

	rows, err := func() (rows [][]string, err error) {
		rowsScan, err := e.excelFile.Rows(config.SheetName)
		if err != nil {
			return nil, fmt.Errorf("excelFile.Rows: %w", err)
		}
		defer func() {
			if closeErr := rowsScan.Close(); closeErr != nil && err == nil {
				err = fmt.Errorf("close: %w", closeErr)
			}
		}()

		var _i = 0
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
				return nil, fmt.Errorf("columns: %w", err1)
			}

			rows = append(rows, row)
		}
		return rows, nil
	}()
	if err != nil {
		return nil, err
	}

	if len(config.ColumnIndexes) > 0 {
		var newRows [][]string
		for _, row := range rows {
			rowLen := len(row)
			if rowLen <= 0 {
				continue
			}

			var newRow []string
			for _, index := range config.ColumnIndexes {
				if index >= rowLen {
					newRow = append(newRow, "")
				} else {
					newRow = append(newRow, row[index])
				}
			}

			newRows = append(newRows, newRow)
		}

		rows = newRows
	}

	return rows, nil
}

// Close releases the underlying excelize file handle.
func (e *Excel) Close() error {
	if err := e.excelFile.Close(); err != nil {
		return fmt.Errorf("excelFile.Close: %w", err)
	}

	return nil
}
