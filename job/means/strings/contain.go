package strings

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/job/means"
	"github.com/auho/go-etl/v2/tool/maps"
)

var _ means.InsertMeans = (*Contains)(nil)
var _ means.UpdateMeans = (*Contains)(nil)

type Contains struct {
	rule means.Ruler

	keyName       string
	items         []map[string]string
	keys          []string
	defaultValues map[string]any
}

func NewContains(rule means.Ruler) *Contains {
	return &Contains{rule: rule}
}

func (c *Contains) GetTitle() string {
	return fmt.Sprintf("Tag:%s{%s}", c.rule.Name(), strings.Join(c.rule.Labels(), ", "))
}

func (c *Contains) GetKeys() []string {
	return c.keys
}

func (c *Contains) DefaultValues() map[string]any {
	return c.defaultValues
}

func (c *Contains) Insert(contents []string) []map[string]any {
	var results []map[string]any

	for _, content := range contents {
		for _, item := range c.items {
			if strings.Contains(content, item[c.keyName]) {
				results = append(results, maps.MapToMapAny(item))
			}
		}
	}

	return results
}

func (c *Contains) Update(contents []string) map[string]any {
	var result map[string]any

	for _, content := range contents {
		for _, item := range c.items {
			if strings.Contains(content, item[c.keyName]) {
				result = maps.MapToMapAny(item)
				break
			}
		}
	}

	return result
}

func (c *Contains) Prepare() error {
	var err error
	c.keyName = c.rule.KeywordNameAlias()
	c.items, err = c.rule.ItemsAlias()
	if err != nil {
		return fmt.Errorf("contains Prepare error; %w", err)
	}

	c.keys = c.rule.MeansKeys()
	c.defaultValues = c.rule.MeansDefaultValues()

	return nil
}

func (c *Contains) Close() error { return nil }
