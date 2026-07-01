package regexps

import (
	"fmt"
	"regexp"

	"github.com/auho/go-etl/v3/job/extract"
)

var _ extract.Extractor = (*SubMatch)(nil)

// regexp sub match
//
// - ab.*cd => 匹配 ab.*cd 部分
// - a(b.*c)d => 匹配 b.*c 部分
//

// SubMatch
// sub match
type SubMatch struct {
	expressions []string
	rule        extract.Rule
	subMode     func([]*regexp.Regexp, []string) results
	toMaps      func(results, extract.Rule) []map[string]any
	keys        []string
	defaults    map[string]any

	regexps []*regexp.Regexp
}

func NewSubMatch(exs []string, rule extract.Rule, subMode func([]*regexp.Regexp, []string) results, toMaps func(results, extract.Rule) []map[string]any, keys []string, defaults map[string]any) *SubMatch {
	return &SubMatch{
		expressions: exs,
		rule:        rule,
		subMode:     subMode,
		toMaps:      toMaps,
		keys:        keys,
		defaults:    defaults,
	}
}

func (r *SubMatch) Title() string {
	return fmt.Sprintf("SubMatch[%s]", r.rule.Name())
}

func (r *SubMatch) Prepare() error {
	for _, ex := range r.expressions {
		r.regexps = append(r.regexps, regexp.MustCompile(ex))
	}

	return nil
}

func (r *SubMatch) Keys() []string {
	return r.keys
}

func (r *SubMatch) DefaultValues() map[string]any {
	return r.defaults
}

func (r *SubMatch) Extract(contents []string) extract.Result {
	defer func() {
		if v := recover(); v != nil {
			panic(fmt.Errorf("do[%#v]", r.expressions))
		}
	}()

	rets := r.subMode(r.regexps, contents)
	if len(rets) == 0 {
		return extract.Result{}
	}
	return extract.NewResult(true, r.toMaps(rets, r.rule))
}

func (r *SubMatch) Close() error { return nil }

func _subMatch(ret []string) (result, bool) {
	var has bool
	var text string

	retLen := len(ret)
	if retLen == 1 {
		text = ret[0]
		has = true
	} else if retLen > 1 {
		text = ret[1]
		has = true
	}

	var result result
	if has {
		result.text = text
		result.amount = 1
	}

	return result, has
}

func _mergeResults(rs results) results {
	if rs == nil {
		return nil
	}

	var newResults results
	resultFlag := make(map[string]int)

	for _, ret := range rs {
		if index, ok := resultFlag[ret.text]; ok {
			newResults[index].amount += 1
		} else {
			newResults = append(newResults, ret)
			resultFlag[ret.text] = len(newResults) - 1
		}
	}

	return newResults
}

// allSubMode finds all sub matches of all contents.
func allSubMode(regexps []*regexp.Regexp, contents []string) results {
	var rets results
	for _, content := range contents {
		for _, re := range regexps {
			ret := re.FindAllStringSubmatch(content, -1)
			if ret != nil {
				for _, _ret := range ret {
					if result, ok := _subMatch(_ret); ok {
						rets = append(rets, result)
					}
				}
			}
		}
	}

	return _mergeResults(rets)
}

// leftmostSubMode finds the leftmost sub match of all contents.
func leftmostSubMode(regexps []*regexp.Regexp, contents []string) results {
	var rets results
	for _, content := range contents {
		for _, re := range regexps {
			ret := re.FindStringSubmatch(content)
			if result, ok := _subMatch(ret); ok {
				rets = append(rets, result)
			}
		}
	}

	return _mergeResults(rets)
}

// firstSubMode finds the leftmost sub match of the first matching content.
func firstSubMode(regexps []*regexp.Regexp, contents []string) results {
	var rets results
	for _, content := range contents {
		for _, re := range regexps {
			ret := re.FindStringSubmatch(content)
			if result, ok := _subMatch(ret); ok {
				rets = append(rets, result)
				goto LOOP
			}
		}
	}
LOOP:

	return _mergeResults(rets)
}

// NewAllSubMatch
// all sub match of all contents
func NewAllSubMatch(exs []string, rule extract.Rule) *SubMatch {
	return NewSubMatch(exs, rule, allSubMode,
		func(r results, rule extract.Rule) []map[string]any { return r.toAll(rule) },
		[]string{rule.NameAlias(), rule.KeywordAmountNameAlias()},
		map[string]any{rule.NameAlias(): "", rule.KeywordAmountNameAlias(): 0})
}

// NewSubMatchAll
// leftmost sub match of all contents
func NewSubMatchAll(exs []string, rule extract.Rule) *SubMatch {
	return NewSubMatch(exs, rule, leftmostSubMode,
		func(r results, rule extract.Rule) []map[string]any { return r.toAll(rule) },
		[]string{rule.NameAlias(), rule.KeywordAmountNameAlias()},
		map[string]any{rule.NameAlias(): "", rule.KeywordAmountNameAlias(): 0})
}

// NewSubMatchAllLine
// leftmost sub match of all contents, line export
func NewSubMatchAllLine(exs []string, rule extract.Rule) *SubMatch {
	return NewSubMatch(exs, rule, leftmostSubMode,
		func(r results, rule extract.Rule) []map[string]any { return r.toLine(rule) },
		[]string{rule.NameAlias(), rule.KeywordNumNameAlias(), rule.KeywordAmountNameAlias()},
		map[string]any{rule.NameAlias(): "", rule.KeywordNumNameAlias(): 0, rule.KeywordAmountNameAlias(): 0})
}

// NewSubMatchAllFlag
// leftmost sub match of all contents, flag export
func NewSubMatchAllFlag(exs []string, rule extract.Rule) *SubMatch {
	return NewSubMatch(exs, rule, leftmostSubMode,
		func(r results, rule extract.Rule) []map[string]any { return r.toFlag(rule) },
		[]string{rule.NameAlias(), rule.KeywordNameAlias()},
		map[string]any{rule.NameAlias(): 0, rule.KeywordNameAlias(): ""})
}

// NewSubMatchFirst
// leftmost sub match of first match found content
func NewSubMatchFirst(exs []string, rule extract.Rule) *SubMatch {
	return NewSubMatch(exs, rule, firstSubMode,
		func(r results, rule extract.Rule) []map[string]any { return r.toAll(rule) },
		[]string{rule.NameAlias(), rule.KeywordAmountNameAlias()},
		map[string]any{rule.NameAlias(): "", rule.KeywordAmountNameAlias(): 0})
}

// NewSubMatchFirstFlag
// leftmost sub match of first match found content, flag export
func NewSubMatchFirstFlag(exs []string, rule extract.Rule) *SubMatch {
	return NewSubMatch(exs, rule, firstSubMode,
		func(r results, rule extract.Rule) []map[string]any { return r.toFlag(rule) },
		[]string{rule.NameAlias(), rule.KeywordNameAlias()},
		map[string]any{rule.NameAlias(): 0, rule.KeywordNameAlias(): ""})
}
