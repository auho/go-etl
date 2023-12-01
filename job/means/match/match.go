package match

import (
	"fmt"
	"sort"
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

func NewKey(rule means.Ruler) *Match {
	return NewMatch(rule, func(rule means.Ruler, m *matcher, c []string) []map[string]any {
		res := m.MatchKey(c)
		if res == nil {
			return nil
		}

		return res.toSliceMapAny(rule)
	})
}

func NewFirstKey(rule means.Ruler) *Match {
	return NewMatch(rule, func(rule means.Ruler, m *matcher, c []string) []map[string]any {
		res := m.MatchFirstKey(c)
		if res == nil {
			return nil
		}

		return res.toSliceMapAny(rule)
	})
}

// NewWholeLabels
// merge all labels together
// label1|label2|label3
// keyword1|keyword2|keyword3|
func NewWholeLabels(rule means.Ruler) *Match {
	return NewMatch(rule, func(rule means.Ruler, m *matcher, c []string) []map[string]any {
		_res := m.MatchLabel(c)
		if _res == nil {
			return nil
		}

		sort.SliceStable(_res, func(i, j int) bool {
			return _res[i].Identity < _res[j].Identity
		})

		_rts := make(map[string][]string)
		_labelAmount := 0
		_keywordAmount := 0

		for _, _r := range _res {
			for _labelKey, _labelValue := range _r.Labels {
				_rts[_labelKey] = append(_rts[_labelKey], _labelValue)
			}

			for _, _key := range _r.Keys {
				_rts[rule.KeywordNameAlias()] = append(_rts[rule.KeywordNameAlias()], _key)

				_keywordAmount += 1
			}

			_labelAmount += 1
		}

		_rt := make(map[string]any)
		for _rk, _rv := range _rts {
			_rt[_rk] = strings.Join(_rv, "|")
		}

		_rt[rule.LabelNumNameAlias()] = _labelAmount
		_rt[rule.KeywordNumNameAlias()] = _keywordAmount

		return []map[string]any{_rt}
	})
}
