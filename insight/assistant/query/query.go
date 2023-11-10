package query

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant/excel/write"
	"github.com/auho/go-etl/v2/insight/assistant/query/dataset"
	"github.com/auho/go-etl/v2/insight/assistant/query/source"
	"github.com/auho/go-toolkit/time/timing"
)

type subQuery struct {
	source      source.Sourcer
	datasetMode dataset.Mode
	state       sqlState
}

type Query struct {
	xlsxPath string
	xlsxName string
	excel    *write.Excel

	queries  []subQuery
	state    state
	duration *timing.Duration
}

func NewQuery(xlsxName, xlsxPath string) (*Query, error) {
	q := &Query{}
	q.xlsxName = xlsxName + ".xlsx"
	q.xlsxPath = xlsxPath

	q.duration = timing.NewDuration()
	q.duration.Start()

	var err error
	q.excel, err = write.NewExcel(fmt.Sprintf("%s/%s", q.xlsxPath, q.xlsxName))
	if err != nil {
		return nil, fmt.Errorf("NewExcel error; %w", err)
	}

	return q, nil
}

// AddAppend
// add append dataset
func (q *Query) AddAppend(source source.Sourcer) {
	q.add(dataset.ModeAppend, source)
}

// AddSpread
// add spread dataset
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
	fmt.Println(fmt.Sprintf("%s/%s", q.xlsxPath, q.xlsxName))
	fmt.Println()

	for _, sq := range q.queries {
		err := q.doQuery(sq)
		if err != nil {
			return fmt.Errorf("doQuery error; %w", err)
		}
	}

	return nil
}

func (q *Query) doQuery(sq subQuery) error {
	_d := timing.NewDuration()
	_d.Start()

	_d.Begin()
	ds, err := sq.source.Dataset()
	sq.state.sourceDuration = _d.SubBegin()
	if err != nil {
		return fmt.Errorf("dataset error; %w", err)
	}

	_d.Begin()
	dsMode, err := dataset.NewMode(sq.datasetMode, ds)
	if err != nil {
		return fmt.Errorf("NewMode error; %w", err)
	}

	data, err := dsMode.Data()
	sq.state.datasetDuration = _d.SubBegin()
	if err != nil {
		return fmt.Errorf("data error; %w", err)
	}

	_d.Begin()
	for _, name := range data.Names {
		_, err = q.excel.NewSheetWithData(name, data.Rows[name])
		sq.state.toSheetDuration = _d.SubBegin()
		if err != nil {
			return fmt.Errorf("NewSheetWithData error; %w", err)
		}
	}

	_d.Stop()
	sq.state.totalDuration = _d.SubStart()
	q.state.add(sq.state)

	fmt.Println(fmt.Sprintf(" 《%s》; %s", dsMode.Name(), sq.state.overview()))
	for _, set := range dsMode.Sets() {
		fmt.Println(fmt.Sprintf("    <%s> => amount: %d, duration %s", set.ItemName, set.Amount, timing.PrettyDuration(set.Duration)))
		fmt.Println(fmt.Sprintf("    %s", set.Sql))
		fmt.Println()
	}

	return nil
}

func (q *Query) Save() error {
	q.duration.Begin()
	err := q.doQueries()
	q.state.queriesDuration = q.duration.SubBegin()
	if err != nil {
		return fmt.Errorf("doQueries error; %w", err)
	}

	q.duration.Begin()
	err = q.excel.SaveAs()
	q.state.saveDuration = q.duration.SubBegin()
	q.duration.Stop()
	q.state.totalDuration = q.duration.SubStart()

	fmt.Println("QUERY =>")
	fmt.Println(q.state.overview())

	return err
}
