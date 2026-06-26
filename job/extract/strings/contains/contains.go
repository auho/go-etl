package contains

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v3/job/extract"
)

var _ extract.Extractor = (*Contains)(nil)

type Contains struct {
	subs   []string
	export *extract.Exporter[Results]

	subMode func([]string) Results
}

func newContains(subs []string, subMode func([]string) Results, export *extract.Exporter[Results]) *Contains {
	return &Contains{
		subs:    subs,
		subMode: subMode,
		export:  export,
	}
}

func (c *Contains) Prepare() error { return nil }

func (c *Contains) Title() string {
	return fmt.Sprintf("Contains[%s]", c.export.GetRule().Name())
}

func (c *Contains) NewExport() extract.FieldSpec {
	return c.export
}

func (c *Contains) Search(contents []string) extract.Result {
	rets := c.subMode(contents)

	return c.export.ToToken(rets, len(rets) > 0)
}

func (c *Contains) Close() error { return nil }

// NewContainsAll
// all sub of all contents
func NewContainsAll(subs []string, export *extract.Exporter[Results]) *Contains {
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
func NewContainsFirst(subs []string, export *extract.Exporter[Results]) *Contains {
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
