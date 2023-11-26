package match

import (
	"strings"

	"github.com/auho/go-etl/v2/tool/maps"
)

type matcher struct {
	keyName string
	items   []map[string]string
}

func newMatcher(keyName string, items []map[string]string) *matcher {
	return &matcher{keyName: keyName, items: items}
}

func (m *matcher) findAll(contents []string) []map[string]any {
	var results []map[string]any
	for _, content := range contents {
		for _, item := range m.items {
			if strings.Contains(content, item[m.keyName]) {
				results = append(results, maps.MapToMapAny(item))
			}
		}
	}

	return results
}

func (m *matcher) findFirst(contents []string) []map[string]any {
	var results []map[string]any
	for _, content := range contents {
		for _, item := range m.items {
			if strings.Contains(content, item[m.keyName]) {
				results = append(results, maps.MapToMapAny(item))
				break
			}
		}
	}

	return results
}
