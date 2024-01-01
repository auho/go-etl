package contains

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
)

var _ search.Searcher = (*Contains)(nil)
var _ search.Exporter = (*Export)(nil)

type Contains struct {
	rule   means.Ruler
	subs   []string
	export *Export

	subMode func([]string) Results
}

func newContains(rule means.Ruler, subs []string, subMode func([]string) Results, export *Export) *Contains {
	return &Contains{
		rule:    rule,
		subs:    subs,
		subMode: subMode,
		export:  export,
	}
}

func NewContainsAll(rule means.Ruler, subs []string, export *Export) *Contains {
	return newContains(rule, subs, func(contents []string) Results {
		var results Results
		for _, content := range contents {
			for _, sub := range subs {
				_c := strings.Count(content, sub)
				if _c > 0 {
					results = append(results, Result{
						Sub:    sub,
						Amount: _c,
					})
				}
			}
		}

		var newResults Results
		resultFlag := make(map[string]int)

		for _, result := range results {
			if index, ok := resultFlag[result.Sub]; ok {
				newResults[index].Amount += 1
			} else {
				newResults = append(newResults, result)
				resultFlag[result.Sub] = len(newResults) - 1
			}
		}

		return newResults
	}, export)
}

func NewContainsAny(rule means.Ruler, subs []string, export *Export) *Contains {
	return newContains(rule, subs, func(contents []string) Results {
		var results Results
		for _, content := range contents {
			for _, sub := range subs {
				_c := strings.Count(content, sub)
				if _c > 0 {
					results = append(results, Result{
						Sub:    sub,
						Amount: _c,
					})

					break
				}
			}
		}

		return results
	}, export)
}

func (c *Contains) Prepare() error { return nil }

func (c *Contains) GetTitle() string {
	return fmt.Sprintf("Contains[%s]", c.rule.Name())
}

func (c *Contains) GenExport() search.Exporter {
	return c.export
}

func (c *Contains) Do(contents []string) search.Token {
	rets := c.subMode(contents)

	return c.export.ToToken(c.rule, rets)
}

func (c *Contains) Close() error { return nil }

func (c *Contains) ToMeans() *means.Means {
	return means.NewMeans(c)
}
