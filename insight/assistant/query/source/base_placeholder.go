package source

import (
	"fmt"
	"regexp"
	"strings"
)

type basePlaceHolder struct {
}

func (bph *basePlaceHolder) buildKeys(sql string) []string {
	re := regexp.MustCompile(`##[^#]+##`)
	keys := re.FindAllString(sql, -1)
	for i, key := range keys {
		keys[i] = strings.ReplaceAll(key, "#", "")
	}

	return keys
}

// []string []item id
// map[string]string map[item id]sql
func (bph *basePlaceHolder) buildPlaceholderItemsSqlSet(s Source, sql string, keys []string, items []map[string]any) ([]string, map[string]string) {
	var itemsId []string
	itemsSql := make(map[string]string, len(items))

	// remove duplicates
	_itemsIdMap := make(map[string]struct{})

	for _, item := range items {
		var itemValues []string
		for _, key := range keys {
			itemValues = append(itemValues, fmt.Sprintf("%v", item[key]))
		}

		itemId := s.itemValuesToIdentification(itemValues)
		if _, ok := _itemsIdMap[itemId]; ok {
			continue
		}

		_itemsIdMap[itemId] = struct{}{}

		itemsId = append(itemsId, itemId)
		itemsSql[itemId] = bph.buildPlaceholderItemSql(sql, item)
	}

	return itemsId, itemsSql
}

// []string sql
func (bph *basePlaceHolder) buildPlaceholderItemsSqlList(sql string, items []map[string]any) []string {
	var itemsSql []string
	for _, item := range items {
		itemsSql = append(itemsSql, bph.buildPlaceholderItemSql(sql, item))
	}

	return itemsSql
}

// string sql
func (bph *basePlaceHolder) buildPlaceholderItemSql(sql string, item map[string]any) string {
	for key, value := range item {
		sql = strings.ReplaceAll(sql, fmt.Sprintf("##%s##", key), fmt.Sprintf("%v", value))
	}

	return sql
}
