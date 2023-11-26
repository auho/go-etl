package match

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/job/means"
)

var _ means.InsertMeans = (*Match)(nil)
var _ means.UpdateMeans = (*Match)(nil)

type Match struct {
	rule    means.Ruler
	matcher *matcher
	fn      func(means.Ruler, *matcher, []string) []map[string]any

	keys          []string
	defaultValues map[string]any
}

func NewMatch(rule means.Ruler, fn func(means.Ruler, *matcher, []string) []map[string]any) *Match {
	return &Match{rule: rule, fn: fn}
}

func (c *Match) GetTitle() string {
	return fmt.Sprintf("Tag:%s{%s}", c.rule.Name(), strings.Join(c.rule.Labels(), ", "))
}

func (c *Match) GetKeys() []string {
	return c.keys
}

func (c *Match) DefaultValues() map[string]any {
	return c.defaultValues
}

func (c *Match) Insert(contents []string) []map[string]any {
	return c.fn(c.rule, c.matcher, contents)
}

func (c *Match) Update(contents []string) map[string]any {
	rs := c.fn(c.rule, c.matcher, contents)
	if rs == nil {
		return nil
	}

	return rs[0]
}

func (c *Match) Prepare() error {
	items, err := c.rule.ItemsAlias()
	if err != nil {
		return fmt.Errorf("contains Prepare error; %w", err)
	}

	c.matcher = newMatcher(c.rule.KeywordNameAlias(), items)

	c.keys = c.rule.MeansKeys()
	c.defaultValues = c.rule.MeansDefaultValues()

	return nil
}

func (c *Match) Close() error { return nil }

func NewFindAll(rule means.Ruler) *Match {
	return NewMatch(rule, func(rule means.Ruler, m *matcher, c []string) []map[string]any {
		return m.findAll(c)
	})
}

func NewFindFirst(rule means.Ruler) *Match {
	return NewMatch(rule, func(rule means.Ruler, m *matcher, c []string) []map[string]any {
		return m.findFirst(c)
	})
}
