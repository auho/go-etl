package query

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/insight/model/tool/maps"
	"github.com/auho/go-etl/v2/insight/model/tool/slices"
	simpleDb "github.com/auho/go-simple-db/v2"
)

type sourcer interface {
	GetSheetName() string
	Rows() ([][]any, error)
}

type Source struct {
	SheetName string
	DB        *simpleDb.SimpleDB
}

func (s *Source) GetSheetName() string {
	return s.SheetName
}

func (s *Source) buildPlaceholderItemsSqlList(sql string, items []map[string]string) []string {
	var itemsSql []string
	for _, item := range items {
		itemsSql = append(itemsSql, s.buildPlaceholderItemSql(sql, item))
	}

	return itemsSql
}

func (s *Source) buildPlaceholderItemSql(sql string, item map[string]string) string {
	for key, value := range item {
		sql = strings.ReplaceAll(sql, fmt.Sprintf("##%s##", key), value)
	}

	return sql
}

func (s *Source) querySql(sql string, fields []string) ([][]any, error) {
	var rows []map[string]any
	err := s.DB.Raw(sql).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("raw error; %w", err)
	}

	return maps.SliceMapStringAnyToSliceSliceAny(rows, fields), nil
}

func (s *Source) rowsAppend(sqlList []string, fields []string) ([][]any, error) {
	var rows [][]any
	rows = append(rows, slices.SliceToAny(fields))

	for _, sql := range sqlList {
		sRows, err := s.querySql(sql, fields)
		if err != nil {
			return nil, fmt.Errorf("querySql error; %w", err)
		}

		rows = append(rows, sRows...)
	}

	return rows, nil
}
