package source

import (
	"fmt"
	"regexp"
	"strings"
)

type basePlaceHolder struct {
}

// []string []item name
// map[string]string map[item name]sql
func (bph *basePlaceHolder) buildPlaceholderItemsSqlSet(s *Source, sql string, items []map[string]string) ([]string, map[string]string) {
	var itemsName []string
	itemsSql := make(map[string]string, len(items))

	re := regexp.MustCompile(`##[^#]+##`)
	keys := re.FindAllString(sql, -1)
	for i, key := range keys {
		keys[i] = strings.ReplaceAll(key, "#", "")
	}

	for _, item := range items {
		var itemsKey []string
		for _, key := range keys {
			itemsKey = append(itemsKey, item[key])
		}

		itemNameId := s.itemsNameToIdentification(itemsKey)
		itemsName = append(itemsName, itemNameId)
		itemsSql[itemNameId] = bph.buildPlaceholderItemSql(sql, item)
	}

	return itemsName, itemsSql
}

// []string sql
func (bph *basePlaceHolder) buildPlaceholderItemsSqlList(sql string, items []map[string]string) []string {
	var itemsSql []string
	for _, item := range items {
		itemsSql = append(itemsSql, bph.buildPlaceholderItemSql(sql, item))
	}

	return itemsSql
}

// string sql
func (bph *basePlaceHolder) buildPlaceholderItemSql(sql string, item map[string]string) string {
	for key, value := range item {
		sql = strings.ReplaceAll(sql, fmt.Sprintf("##%s##", key), value)
	}

	return sql
}
