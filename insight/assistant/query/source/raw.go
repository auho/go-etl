package source

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/query/dataset"
)

type RawSource struct {
	Name string
	Raw  assistant.Rawer

	source Source
}

func NewRaw(name string, raw assistant.Rawer) *RawSource {
	return &RawSource{
		Name: name,
		Raw:  raw,
		source: Source{
			Name: name,
			DB:   raw.GetDB(),
		},
	}
}

func (rs *RawSource) Dataset() (*dataset.Dataset, error) {
	fields, err := rs.Raw.GetDB().GetTableColumns(rs.Raw.TableName())
	if err != nil {
		return nil, err
	}

	itemsId := []string{rs.Name}
	itemsSql := map[string]string{rs.Name: rs.Raw.DmlTable().Select([]string{"*"}).Sql()}

	sets, err := rs.source.queryItemsSet(
		fields,
		itemsId,
		itemsSql,
	)
	if err != nil {
		return nil, fmt.Errorf("queryItemsSet error; %w", err)
	}

	return &dataset.Dataset{
		Name:   rs.Name,
		Titles: fields,
		Sets:   sets,
	}, nil
}
