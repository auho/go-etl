package regexps

import (
	"fmt"
	"regexp"

	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
)

var _ search.Searcher = (*SubMatch)(nil)

type SubMatch struct {
	rule        means.Ruler
	expressions []string
	subMode     func([]*regexp.Regexp, []string) Results
	export      *Export

	regexps []*regexp.Regexp
}

func NewSubMatch(rule means.Ruler, exs []string, subMode func([]*regexp.Regexp, []string) Results, export *Export) *SubMatch {
	return &SubMatch{
		rule:        rule,
		expressions: exs,
		subMode:     subMode,
		export:      export,
	}
}

func (r *SubMatch) GetTitle() string {
	return "SubMatch"
}

func (r *SubMatch) Prepare() error {
	for _, ex := range r.expressions {
		r.regexps = append(r.regexps, regexp.MustCompile(ex))
	}

	return nil
}

func (r *SubMatch) GenExport() search.Exporter {
	return r.export
}

func (r *SubMatch) Do(contents []string) search.Token {
	defer func() {
		if v := recover(); v != nil {
			panic(fmt.Errorf("SubMatch expressions[%#v]", r.expressions))
		}
	}()

	rets := r.subMode(r.regexps, contents)

	return r.export.ToToken(rets, r.rule)
}

func (r *SubMatch) Close() error { return nil }

func (r *SubMatch) ToMeans() *means.Means {
	return means.NewMeans(r)
}

func _subMatch(ret []string) (Result, bool) {
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

	var result Result
	if has {
		result.Text = text
		result.Amount = 1
	}

	return result, has
}

func _mergeResults(results Results) Results {
	if results == nil {
		return nil
	}

	var newResults Results
	resultFlag := make(map[string]int)

	for _, result := range results {
		if index, ok := resultFlag[result.Text]; ok {
			newResults[index].Amount += 1
		} else {
			newResults = append(newResults, result)
			resultFlag[result.Text] = len(newResults) - 1
		}
	}

	return newResults
}

func NewAllSubMatch(rule means.Ruler, exs []string, export *Export) *SubMatch {
	return NewSubMatch(rule, exs, func(regexps []*regexp.Regexp, contents []string) Results {
		var rets Results
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
	}, export)
}

func NewSubMatchAll(rule means.Ruler, exs []string, export *Export) *SubMatch {
	return NewSubMatch(rule, exs, func(regexps []*regexp.Regexp, contents []string) Results {
		var rets Results
		for _, content := range contents {
			for _, re := range regexps {
				ret := re.FindStringSubmatch(content)
				if result, ok := _subMatch(ret); ok {
					rets = append(rets, result)
				}
			}
		}

		return _mergeResults(rets)
	}, export)
}

func NewSubMatchFirst(rule means.Ruler, exs []string, export *Export) *SubMatch {
	return NewSubMatch(rule, exs, func(regexps []*regexp.Regexp, contents []string) Results {
		var rets Results
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
	}, export)
}
