package source

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/auho/go-etl/v2/insight/assistant/query/dataset"
	"github.com/auho/go-etl/v2/tool/maps"
	simpleDb "github.com/auho/go-simple-db/v2"
)

type Sourcer interface {
	Dataset() (*dataset.Dataset, error)
}

type Source struct {
	HasNamePrefix bool
	Name          string
	DB            *simpleDb.SimpleDB
}

func (s *Source) buildPlaceholderItemsSqlSet(sql string, items []map[string]string) ([]string, map[string]string) {
	var itemsName []string
	itemsSql := make(map[string]string, len(items))

	re := regexp.MustCompile(`##[^#]+##`)
	keys := re.FindAllString(sql, -1)
	for i, key := range keys {
		keys[i] = strings.ReplaceAll(key, "#", "")
	}

	for _, item := range items {
		var itemKey []string
		for _, key := range keys {
			itemKey = append(itemKey, item[key])
		}

		itemName := strings.Join(itemKey, "_")
		if s.HasNamePrefix {
			itemName = fmt.Sprintf("%s_%s", s.Name, itemName)
		}

		itemsName = append(itemsName, itemName)
		itemsSql[itemName] = s.buildPlaceholderItemSql(sql, item)
	}

	return itemsName, itemsSql
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

func (s *Source) queryItemsSet(fields, itemsName []string, itemsSql map[string]string) (map[string][][]any, error) {
	m := make(map[string][][]any, len(itemsName))

	for _, itemName := range itemsName {
		rows, err := s.querySql(itemsSql[itemName], fields)
		if err != nil {
			return nil, fmt.Errorf("querySql error; %w", err)
		}

		m[itemName] = rows
	}

	return m, nil
}

func (s *Source) querySql(sql string, fields []string) ([][]any, error) {
	var rows []map[string]any
	err := s.DB.Raw(sql).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("raw error; %w", err)
	}

	return maps.SliceMapStringAnyToSliceSliceAny(rows, fields), nil
}
