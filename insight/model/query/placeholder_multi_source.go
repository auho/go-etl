package query

import (
	"github.com/auho/go-etl/v2/insight/model/dml"
)

var _ sourcer = (*PlaceholderListSource)(nil)

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

func (pls *PlaceholderListSource) Rows() ([][]any, error) {
	fields := pls.Table.GetSelectFields()
	sql := pls.Table.Sql()

	sqls := pls.buildPlaceholderItemsSqlList(sql, pls.Items)
	return pls.rowsAppend(sqls, fields)
}
