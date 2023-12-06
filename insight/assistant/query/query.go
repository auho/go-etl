package query

import (
	"fmt"
	"path"

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
	excel    *write.Excel

	queries  []*subQuery
	state    state
	duration *timing.Duration

	summary []string
}

func Do(xlsxName, xlsxPath string, fn func(q *Query)) error {
	return DoWithPath(path.Join(xlsxPath, xlsxName+".xlsx"), fn)
}

func DoWithPath(xlsxFilePath string, fn func(q *Query)) error {
	q, err := NewQueryWithPath(xlsxFilePath)
	if err != nil {
		return err
	}

	fn(q)

	err = q.Save()
	if err != nil {
		return err
	}

	return nil
}

func NewQueryWithPath(xlsxFilePath string) (*Query, error) {
	q := &Query{}
	q.xlsxPath = xlsxFilePath

	q.duration = timing.NewDuration()
	q.duration.Start()

	var err error
	q.excel, err = write.NewExcel(q.xlsxPath)
	if err != nil {
		return nil, fmt.Errorf("NewExcel error; %w", err)
	}

	return q, nil
}

// NewQuery
// xlsx Name: without file suffix
func NewQuery(xlsxName, xlsxPath string) (*Query, error) {
	return NewQueryWithPath(path.Join(xlsxPath, xlsxName+".xlsx"))
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
	q.queries = append(q.queries, &subQuery{
		source:      s,
		datasetMode: dm,
	})
}

func (q *Query) doQueries() error {
	fmt.Println(q.xlsxPath)
	fmt.Println()

	for _, sq := range q.queries {
		err := q.doQuery(sq)
		if err != nil {
			return fmt.Errorf("doQuery error; %w", err)
		}
	}

	return nil
}

func (q *Query) doQuery(sq *subQuery) error {
	_d := timing.NewDuration()
	_d.Start()

	_d.Begin()
	_dataset, err := sq.source.Dataset()
	sq.state.sourceDuration = _d.SubBegin()
	if err != nil {
		return fmt.Errorf("dataset error; %w", err)
	}

	_d.Begin()
	_datasetMode, err := dataset.NewMode(sq.datasetMode, _dataset)
	if err != nil {
		return fmt.Errorf("NewMode error; %w", err)
	}

	_data, err := _datasetMode.Data()
	sq.state.datasetDuration = _d.SubBegin()
	if err != nil {
		return fmt.Errorf("data error; %w", err)
	}

	_d.Begin()
	for _, name := range _data.Names {
		_, err = q.excel.NewSheetWithData(name, _data.Rows[name])
		sq.state.toSheetDuration = _d.SubBegin()
		if err != nil {
			return fmt.Errorf("NewSheetWithData error; %w", err)
		}

		sq.state.amount += _data.RowsAmount[name]
	}

	_d.Stop()
	sq.state.totalDuration = _d.SubStart()
	q.state.add(sq.state)

	_querySummary := fmt.Sprintf("《%s》[%s]: %s", _datasetMode.Name(), sq.datasetMode, sq.state.overview())
	q.summary = append(q.summary, _querySummary)
	fmt.Println(_querySummary)

	for _, _set := range _datasetMode.Sets() {
		_setSummary := fmt.Sprintf("  <%s> => amount: %d, duration %s", _set.Name, _set.Amount, timing.PrettyDuration(_set.Duration))
		q.summary = append(q.summary, _setSummary)
		fmt.Println(_setSummary)

		for _, _query := range _set.Queries {
			fmt.Println(fmt.Sprintf("    %s => amount: %d, duration %s:", _query.Name, _query.Amount, timing.PrettyDuration(_query.Duration)))
			fmt.Println(fmt.Sprintf("    %s", _query.Sql))
		}

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

	fmt.Println("SUMMARY =>")
	for _, _s := range q.summary {
		fmt.Println(_s)
	}

	fmt.Println("QUERY =>")
	fmt.Println(q.state.overview())

	return err
}
