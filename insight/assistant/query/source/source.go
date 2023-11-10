package source

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/auho/go-etl/v2/insight/assistant/query/dataset"
	"github.com/auho/go-etl/v2/tool/maps"
	simpleDb "github.com/auho/go-simple-db/v2"
)

type Sourcer interface {
	Dataset() (*dataset.Dataset, error)
}

type Source struct {
	HasNamePrefix bool // Add the name prefix before the item
	Name          string
	DB            *simpleDb.SimpleDB
}

// []string []item name
// map[string]string map[item name]sql
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

// []string sql
func (s *Source) buildPlaceholderItemsSqlList(sql string, items []map[string]string) []string {
	var itemsSql []string
	for _, item := range items {
		itemsSql = append(itemsSql, s.buildPlaceholderItemSql(sql, item))
	}

	return itemsSql
}

// string sql
func (s *Source) buildPlaceholderItemSql(sql string, item map[string]string) string {
	for key, value := range item {
		sql = strings.ReplaceAll(sql, fmt.Sprintf("##%s##", key), value)
	}

	return sql
}

func (s *Source) queryItemsSet(fields, itemsName []string, itemsSql map[string]string) ([]dataset.Set, error) {
	var sets []dataset.Set

	for _, itemName := range itemsName {
		rows, _d, err := s.querySql(itemsSql[itemName], fields)
		if err != nil {
			return nil, fmt.Errorf("querySql error; %w", err)
		}

		sets = append(sets, dataset.Set{
			ItemName: itemName,
			Sql:      itemsSql[itemName],
			Amount:   len(rows),
			Rows:     rows,
			Duration: _d,
		})
	}

	return sets, nil
}

func (s *Source) querySql(sql string, fields []string) ([][]any, time.Duration, error) {
	var rows []map[string]any

	_start := time.Now()
	err := s.DB.Raw(sql).Scan(&rows).Error
	if err != nil {
		return nil, 0, fmt.Errorf("raw error; %w", err)
	}

	return maps.SliceMapStringAnyToSliceSliceAny(rows, fields), time.Now().Sub(_start), nil
}
