package read

import (
	"fmt"
)

type SheetData struct {
	rows [][]string
}

func NewSheetData(excelPath, sheetName string) (*SheetData, error) {
	excel, err := NewExcel(excelPath)
	if err != nil {
		return nil, fmt.Errorf("NewExcel error; %w", err)
	}

	sd := &SheetData{}
	sd.rows, err = excel.excelFile.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("GetRows error; %w", err)
	}

	return sd, nil
}
