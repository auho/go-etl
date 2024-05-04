package contains

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
)

var _ search.Searcher = (*Contains)(nil)

type Contains struct {
	subs   []string
	export *Export

	subMode func([]string) Results
}

func newContains(subs []string, subMode func([]string) Results, export *Export) *Contains {
	return &Contains{
		subs:    subs,
		subMode: subMode,
		export:  export,
	}
}

func (c *Contains) Prepare() error { return nil }

func (c *Contains) GetTitle() string {
	return fmt.Sprintf("Contains[%s]", c.export.GetRule().Name())
}

func (c *Contains) GenExport() search.Exporter {
	return c.export
}

func (c *Contains) Do(contents []string) search.Token {
	rets := c.subMode(contents)

	return c.export.ToToken(rets)
}

func (c *Contains) Close() error { return nil }

func (c *Contains) ToMeans() *means.Means {
	return means.NewMeans(c)
}

// NewContainsAll
// all sub of all contents
func NewContainsAll(subs []string, export *Export) *Contains {
	return newContains(subs, func(contents []string) Results {
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

// NewContainsFirst
// first sub of contents
func NewContainsFirst(subs []string, export *Export) *Contains {
	return newContains(subs, func(contents []string) Results {
		var results Results
		for _, content := range contents {
			for _, sub := range subs {
				_c := strings.Count(content, sub)
				if _c > 0 {
					results = append(results, Result{
						Sub:    sub,
						Amount: _c,
					})

					goto LOOP
				}
			}
		}
	LOOP:
		return results
	}, export)
}
