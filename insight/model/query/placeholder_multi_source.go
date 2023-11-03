package query

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/model/dml"
)

var _ sheets = (*PlaceholderListSource)(nil)

/*
 a: 1 b: 3
 a: 1 b: 3
 a: 2 b: 4
 a: 2 b: 4
*/

type PlaceholderListSource struct {
	Source
	Table dml.Tabler
	Items []map[string]string // []map[field][field value]
}

func (pls *PlaceholderListSource) Sheets() ([]string, map[string][][]any, error) {
	fields := pls.Table.GetSelectFields()
	sql := pls.Table.Sql()

	sqls := pls.buildPlaceholderItemsSqlList(sql, pls.Items)
	rows, err := pls.rowsAppend(sqls, fields)
	if err != nil {
		return nil, nil, fmt.Errorf("rowsAppend error; %w", err)
	}

	return []string{pls.SheetName}, map[string][][]any{pls.SheetName: rows}, nil
}
