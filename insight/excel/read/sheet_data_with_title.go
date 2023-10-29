package read

import (
	"fmt"
)

type SheetDataWithTitle struct {
	titles []string
	rows   [][]string
}

func NewSheetDataWithTitle(excelPath, sheetName string) (*SheetDataWithTitle, error) {
	excel, err := NewExcel(excelPath)
	if err != nil {
		return nil, fmt.Errorf("NewExcel error; %w", err)
	}

	var rows [][]string
	rows, err = excel.excelFile.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("GetRows error; %w", err)
	}

	sd := &SheetDataWithTitle{}
	sd.build(rows)

	return sd, nil
}

func (sdwt *SheetDataWithTitle) build(rows [][]string) {
	sdwt.titles = rows[0]
	titleLen := len(sdwt.titles)

	rows = rows[1:]

	for _, row := range rows {
		if len(row) >= titleLen {
			row = row[0:titleLen]
		}

		sdwt.rows = append(sdwt.rows, row)
	}
}
