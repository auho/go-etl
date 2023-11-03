package query

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/model/excel/write"
	"github.com/auho/go-etl/v2/insight/model/query/dataset"
	"github.com/auho/go-etl/v2/insight/model/query/source"
)

type subQuery struct {
	source      source.Sourcer
	datasetMode dataset.Mode
}

type Query struct {
	xlsxPath string
	xlsxName string
	excel    *write.Excel

	queries []subQuery
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

func (q *Query) AddAppend(source source.Sourcer) {
	q.add(dataset.ModeAppend, source)
}

func (q *Query) AddSpread(source source.Sourcer) {
	q.add(dataset.ModeSpread, source)
}

func (q *Query) add(dm dataset.Mode, s source.Sourcer) {
	q.queries = append(q.queries, subQuery{
		source:      s,
		datasetMode: dm,
	})
}

func (q *Query) doQueries() error {
	for _, sq := range q.queries {
		err := q.doQuery(sq)
		if err != nil {
			return fmt.Errorf("doQuery error; %w", err)
		}
	}

	return nil
}

func (q *Query) doQuery(sq subQuery) error {
	ds, err := sq.source.Dataset()
	if err != nil {
		return fmt.Errorf("dataset error; %w", err)
	}

	dsMode, err := dataset.NewMode(sq.datasetMode, ds)
	if err != nil {
		return fmt.Errorf("NewMode error; %w", err)
	}

	data, err := dsMode.Data()
	if err != nil {
		return fmt.Errorf("data error; %w", err)
	}

	for _, name := range data.Names {
		_, err = q.excel.NewSheetWithData(name, data.Rows[name])
		if err != nil {
			return fmt.Errorf("NewSheetWithData error; %w", err)
		}
	}

	return nil
}

func (q *Query) Save() error {
	err := q.doQueries()
	if err != nil {
		return fmt.Errorf("doQueries error; %w", err)
	}

	return q.excel.SaveAs()
}
