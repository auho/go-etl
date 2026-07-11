package regexps

import (
	"fmt"
	"maps"
	"regexp"

	"github.com/auho/go-etl/v3/task/extract"
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
	toRows      func(results, extract.Rule) []map[string]any
	keys        []string
	defaults    map[string]any

	regexps []*regexp.Regexp
}

func NewSubMatch(exs []string, rule extract.Rule, subMode func([]*regexp.Regexp, []string) results, toRows func(results, extract.Rule) []map[string]any, keys []string, defaults map[string]any) *SubMatch {
	return &SubMatch{
		expressions: exs,
		rule:        rule,
		subMode:     subMode,
		toRows:      toRows,
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
	return maps.Clone(r.defaults)
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
	return extract.NewResult(true, r.toRows(rets, r.rule))
}

func (r *SubMatch) Close() error { return nil }

func subMatch(matched []string) (result, bool) {
	var has bool
	var text string

	retLen := len(matched)
	if retLen == 1 {
		text = matched[0]
		has = true
	} else if retLen > 1 {
		text = matched[1]
		has = true
	}

	var ret result
	if has {
		ret.text = text
		ret.amount = 1
	}

	return ret, has
}

func mergeResults(rets results) results {
	if len(rets) == 0 {
		return nil
	}

	var newRets results
	retFlag := make(map[string]int)

	for _, ret := range rets {
		if index, ok := retFlag[ret.text]; ok {
			newRets[index].amount += 1
		} else {
			newRets = append(newRets, ret)
			retFlag[ret.text] = len(newRets) - 1
		}
	}

	return newRets
}

// allSubMode finds all sub matches of all contents.
func allSubMode(regexps []*regexp.Regexp, contents []string) results {
	var rets results
	for _, content := range contents {
		for _, re := range regexps {
			subMatched := re.FindAllStringSubmatch(content, -1)
			for _, matched := range subMatched {
				if ret, ok := subMatch(matched); ok {
					rets = append(rets, ret)
				}
			}
		}
	}

	return mergeResults(rets)
}

// leftmostSubMode finds the leftmost sub match of all contents.
func leftmostSubMode(regexps []*regexp.Regexp, contents []string) results {
	var rets results
	for _, content := range contents {
		for _, re := range regexps {
			matched := re.FindStringSubmatch(content)
			if ret, ok := subMatch(matched); ok {
				rets = append(rets, ret)
			}
		}
	}

	return mergeResults(rets)
}

// firstSubMode finds the leftmost sub match of the first matching content.
func firstSubMode(regexps []*regexp.Regexp, contents []string) results {
	var rets results
	for _, content := range contents {
		for _, re := range regexps {
			matched := re.FindStringSubmatch(content)
			if ret, ok := subMatch(matched); ok {
				rets = append(rets, ret)
				goto LOOP
			}
		}
	}
LOOP:

	return mergeResults(rets)
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
