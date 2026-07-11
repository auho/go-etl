package source

import (
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant/query/dataset"
)

var _ Source = (*RowsStack)(nil)

// RowsStack
// general stack queries
type RowsStack struct {
	name        string
	rowsSources []*Rows
}

func NewRowsStack(name string, bases ...Base) *RowsStack {
	rs := &RowsStack{}
	rs.name = name

	for _, _b := range bases {
		rs.rowsSources = append(rs.rowsSources, NewRows(_b))
	}

	return rs
}

func (rs *RowsStack) Dataset() (*dataset.Dataset, error) {
	if len(rs.rowsSources) <= 0 {
		return nil, fmt.Errorf("source[%s] rowsSources length is invalid", rs.name)
	}

	var _sets []dataset.Subset

	for _, _rs := range rs.rowsSources {
		ds, err := _rs.Dataset()
		if err != nil {
			return nil, fmt.Errorf("dataset: %w", err)
		}

		_sets = append(_sets, ds.Sets...)
	}

	return &dataset.Dataset{
		Name:   rs.name,
		Titles: rs.rowsSources[0].Table.GetSelectFields(),
		Sets:   _sets,
	}, nil
}
