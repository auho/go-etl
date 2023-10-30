package write

import (
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
	e.excelFile = excelize.NewFile()

	return e, nil
}

func (e *Excel) NewSheetWithData(sheetName string, rows [][]any) (int, error) {
	index, err := e.excelFile.NewSheet(sheetName)
	if err != nil {
		return index, fmt.Errorf("NewSheet error; %w", err)
	}

	var cell string
	for i, row := range rows {
		cell, err = excelize.CoordinatesToCellName(1, i+1) // returns "A1", nil
		if err != nil {
			return -1, fmt.Errorf("CoordinatesToCellName error; %w", err)
		}

		err = e.excelFile.SetSheetRow(sheetName, cell, &row)
		if err != nil {
			return -1, fmt.Errorf("SetSheetRow error; %w", err)
		}
	}

	return -1, nil
}

func (e *Excel) SaveAs() error {
	err := e.excelFile.SaveAs(e.path)
	if err != nil {
		return fmt.Errorf("SaveAs error; %w", err)
	}

	return nil
}

func (e *Excel) Close() error {
	err := e.excelFile.Close()
	if err != nil {
		return fmt.Errorf("close error; %w", err)
	}

	return nil
}
