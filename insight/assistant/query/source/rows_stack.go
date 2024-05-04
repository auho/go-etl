package source

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant/query/dataset"
)

var _ Sourcer = (*RowsStackSource)(nil)

// RowsStackSource
// general stack queries
type RowsStackSource struct {
	name string
	rss  []*RowsSource
}

func NewRowsStack(name string, ss ...Source) *RowsStackSource {
	rs := &RowsStackSource{}
	rs.name = name

	for _, _s := range ss {
		rs.rss = append(rs.rss, NewRows(_s))
	}

	return rs
}

func (rs *RowsStackSource) Dataset() (*dataset.Dataset, error) {
	var _sets []dataset.Set

	for _, _rs := range rs.rss {
		ds, err := _rs.Dataset()
		if err != nil {
			return nil, fmt.Errorf("dataset error; %w", err)
		}

		_sets = append(_sets, ds.Sets...)
	}

	return &dataset.Dataset{
		Name:   rs.name,
		Titles: rs.rss[0].Table.GetSelectFields(),
		Sets:   _sets,
	}, nil
}
