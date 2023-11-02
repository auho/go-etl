package query

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/model/dml"
	"github.com/auho/go-etl/v2/insight/model/tool/maps"
	"github.com/auho/go-etl/v2/insight/model/tool/slices"
)

var _ sourcer = (*TableSource)(nil)

type TableSource struct {
	Source
	Table dml.Tabler
}

func (tq *TableSource) Rows() ([][]any, error) {
	var rows []map[string]any
	err := tq.DB.Raw(tq.Table.Sql()).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("raw error; %w", err)
	}

	var rowsAny [][]any
	fields := tq.Table.GetSelectFields()
	rowsAny = append(rowsAny, slices.SliceToAny(fields))
	rowsAny = append(rowsAny, maps.SliceMapStringAnyToSliceSliceAny(rows, fields)...)

	return rowsAny, nil
}
