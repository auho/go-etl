package source

import (
	"context"
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/query/dataset"
)

type RawSource struct {
	Name string
	Raw  assistant.Raw

	source Base
}

func NewRaw(name string, raw assistant.Raw) *RawSource {
	return &RawSource{
		Name: name,
		Raw:  raw,
		source: Base{
			Name: name,
			DB:   raw.DB(),
		},
	}
}

func (rs *RawSource) Dataset() (*dataset.Dataset, error) {
	fields, err := rs.Raw.DB().GetTableColumns(context.TODO(), rs.Raw.TableName())
	if err != nil {
		return nil, fmt.Errorf("getTableColumns: %w", err)
	}

	itemsId := []string{rs.Name}
	itemsSql := map[string]string{rs.Name: rs.Raw.DMLTable().Select([]string{"*"}).SQL()}

	sets, err := rs.source.queryItemsSet(
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
