package source

import (
	"context"
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/query/dataset"
)

var _ Source = (*RawSource)(nil)

type RawSource struct {
	Base
	Raw assistant.Raw
}

func NewRaw(name string, raw assistant.Raw) *RawSource {
	return &RawSource{
		Base: Base{
			Name: name,
			DB:   raw.DB(),
		},
		Raw: raw,
	}
}

func (rs *RawSource) Dataset() (*dataset.Dataset, error) {
	fields, err := rs.Raw.DB().GetTableColumns(context.TODO(), rs.Raw.TableName())
	if err != nil {
		return nil, fmt.Errorf("getTableColumns: %w", err)
	}

	itemsId := []string{rs.Name}
	itemsSql := map[string]string{rs.Name: rs.Raw.DMLTable().Select([]string{"*"}).SQL()}

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
