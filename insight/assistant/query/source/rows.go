package source

import (
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant/query/dataset"
)

var _ Source = (*Rows)(nil)

// Rows
// general queries
type Rows struct {
	Base
}

func NewRows(b Base) *Rows {
	return &Rows{Base: b}
}

func (rs *Rows) Dataset() (*dataset.Dataset, error) {
	fields := rs.Table.GetSelectFields()
	itemsId := []string{rs.Name}
	itemsSql := map[string]string{rs.Name: rs.Table.SQL()}

	sets, err := rs.queryItemsSet(
		fields,
		itemsId,
		itemsSql,
	)
	if err != nil {
		return nil, fmt.Errorf("queryItemsSet: %w", err)
	}

	return &dataset.Dataset{
		Name:   rs.Name,
		Titles: fields,
		Sets:   sets,
	}, nil
}
