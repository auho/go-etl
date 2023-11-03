package query

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/model/dml"
)

var _ sheets = (*TableSource)(nil)

type TableSource struct {
	Source
	Table dml.Tabler
}

func (ts *TableSource) Sheets() ([]string, map[string][][]any, error) {
	sqls := []string{ts.Table.Sql()}

	rows, err := ts.rowsAppend(sqls, ts.Table.GetSelectFields())
	if err != nil {
		return nil, nil, fmt.Errorf("rowsAppend error; %w", err)
	}

	return []string{ts.SheetName}, map[string][][]any{ts.SheetName: rows}, nil
}
