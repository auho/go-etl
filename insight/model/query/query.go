package query

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/model/excel/write"
)

type Query struct {
	xlsxPath string
	xlsxName string
	excel    *write.Excel

	sources []sourcer
}

func NewQuery(xlsxName, xlsxPath string) (*Query, error) {
	q := &Query{}
	q.xlsxName = xlsxName + ".xlsx"
	q.xlsxPath = xlsxPath

	var err error
	q.excel, err = write.NewExcel(fmt.Sprintf("%s/%s", q.xlsxPath, q.xlsxName))
	if err != nil {
		return nil, fmt.Errorf("NewExcel error; %w", err)
	}

	return q, nil
}

func (q *Query) AddSource(sources ...sourcer) {
	q.sources = append(q.sources, sources...)
}

func (q *Query) doSources() error {
	var err error
	for _, source := range q.sources {
		err = q.doSource(source)
		if err != nil {
			return err
		}
	}

	return nil
}

func (q *Query) doSource(source sourcer) error {
	rows, err := source.Rows()
	if err != nil {
		return fmt.Errorf("rows error; %w", err)
	}

	_, err = q.excel.NewSheetWithData(source.GetSheetName(), rows)
	if err != nil {
		return fmt.Errorf("NewSheetWithData error; %w", err)
	}

	return nil
}

func (q *Query) Save() error {
	err := q.doSources()
	if err != nil {
		return fmt.Errorf("doSources error; %w", err)
	}
	return q.excel.SaveAs()
}
