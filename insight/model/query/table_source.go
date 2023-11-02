package query

import (
	"github.com/auho/go-etl/v2/insight/model/dml"
)

var _ sourcer = (*TableSource)(nil)

type TableSource struct {
	Source
	Table dml.Tabler
}

func (ts *TableSource) Rows() ([][]any, error) {
	sqls := []string{ts.Table.Sql()}

	return ts.rowsAppend(sqls, ts.Table.GetSelectFields())
}
