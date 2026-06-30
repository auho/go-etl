package containment

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v3/job/extract"
)

var _ extract.Extractor = (*Contains)(nil)

type Contains struct {
	subs      []string
	rule      extract.Rule
	subMode   func([]string) results
	toMaps    func(results, extract.Rule) []map[string]any
	keys      []string
	defaults  map[string]any
}

func newContains(subs []string, rule extract.Rule, subMode func([]string) results, toMaps func(results, extract.Rule) []map[string]any, keys []string, defaults map[string]any) *Contains {
	return &Contains{
		subs:     subs,
		rule:     rule,
		subMode:  subMode,
		toMaps:   toMaps,
		keys:     keys,
		defaults: defaults,
	}
}

func (c *Contains) Prepare() error { return nil }

func (c *Contains) Title() string {
	return fmt.Sprintf("Contains[%s]", c.rule.Name())
}

func (c *Contains) Keys() []string {
	return c.keys
}

func (c *Contains) DefaultValues() map[string]any {
	return c.defaults
}

func (c *Contains) Extract(contents []string) extract.Result {
	rets := c.subMode(contents)
	if len(rets) == 0 {
		return extract.Result{}
	}
	return extract.NewResult(true, c.toMaps(rets, c.rule))
}

func (c *Contains) Close() error { return nil }

// allSubMode collects all subs of all contents.
func allSubMode(subs []string) func([]string) results {
	return func(contents []string) results {
		var rs results
		for _, content := range contents {
			for _, sub := range subs {
				_c := strings.Count(content, sub)
				if _c > 0 {
					rs = append(rs, result{
						sub:    sub,
						amount: _c,
					})
				}
			}
		}

		var newResults results
		resultFlag := make(map[string]int)

		for _, result := range rs {
			if index, ok := resultFlag[result.sub]; ok {
				newResults[index].amount += 1
			} else {
				newResults = append(newResults, result)
				resultFlag[result.sub] = len(newResults) - 1
			}
		}

		return newResults
	}
}

// firstSubMode collects the first sub found in contents.
func firstSubMode(subs []string) func([]string) results {
	return func(contents []string) results {
		var rs results
		for _, content := range contents {
			for _, sub := range subs {
				_c := strings.Count(content, sub)
				if _c > 0 {
					rs = append(rs, result{
						sub:    sub,
						amount: _c,
					})

					goto LOOP
				}
			}
		}
	LOOP:
		return rs
	}
}

// NewContainsAll
// all sub of all contents
func NewContainsAll(subs []string, rule extract.Rule) *Contains {
	return newContains(subs, rule, allSubMode(subs),
		func(r results, rule extract.Rule) []map[string]any { return r.toAll(rule) },
		[]string{rule.NameAlias(), rule.KeywordAmountNameAlias()},
		map[string]any{rule.NameAlias(): "", rule.KeywordAmountNameAlias(): 0})
}

// NewContainsAllLine
// all sub of all contents, line export
func NewContainsAllLine(subs []string, rule extract.Rule) *Contains {
	return newContains(subs, rule, allSubMode(subs),
		func(r results, rule extract.Rule) []map[string]any { return r.toLine(rule) },
		[]string{rule.NameAlias(), rule.KeywordNumNameAlias(), rule.KeywordAmountNameAlias()},
		map[string]any{rule.NameAlias(): "", rule.KeywordNumNameAlias(): 0, rule.KeywordAmountNameAlias(): 0})
}

// NewContainsAllFlag
// all sub of all contents, flag export
func NewContainsAllFlag(subs []string, rule extract.Rule) *Contains {
	return newContains(subs, rule, allSubMode(subs),
		func(r results, rule extract.Rule) []map[string]any { return r.toFlag(rule) },
		[]string{rule.NameAlias(), rule.KeywordNameAlias()},
		map[string]any{rule.NameAlias(): 0, rule.KeywordNameAlias(): ""})
}

// NewContainsFirst
// first sub of contents
func NewContainsFirst(subs []string, rule extract.Rule) *Contains {
	return newContains(subs, rule, firstSubMode(subs),
		func(r results, rule extract.Rule) []map[string]any { return r.toAll(rule) },
		[]string{rule.NameAlias(), rule.KeywordAmountNameAlias()},
		map[string]any{rule.NameAlias(): "", rule.KeywordAmountNameAlias(): 0})
}

// NewContainsFirstLine
// first sub of contents, line export
func NewContainsFirstLine(subs []string, rule extract.Rule) *Contains {
	return newContains(subs, rule, firstSubMode(subs),
		func(r results, rule extract.Rule) []map[string]any { return r.toLine(rule) },
		[]string{rule.NameAlias(), rule.KeywordNumNameAlias(), rule.KeywordAmountNameAlias()},
		map[string]any{rule.NameAlias(): "", rule.KeywordNumNameAlias(): 0, rule.KeywordAmountNameAlias(): 0})
}
