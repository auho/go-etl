package read

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

	var err error
	e.excelFile, err = excelize.OpenFile(e.path)
	if err != nil {
		return nil, fmt.Errorf("NewExcel error; %w", err)
	}

	return e, nil
}

func (e *Excel) Close() error {
	if err := e.excelFile.Close(); err != nil {
		return fmt.Errorf("close excel error; %w", err)
	}

	return nil
}
