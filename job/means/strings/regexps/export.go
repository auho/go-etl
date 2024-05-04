package regexps

import (
	"maps"

	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
	maps2 "github.com/auho/go-etl/v2/tool/maps"
)

var _ search.Exporter = (*Export)(nil)

type Export struct {
	rule          means.Ruler
	keys          []string
	defaultValues map[string]any

	resultsToToken func(Results, means.Ruler) []map[string]any
}

func NewExport(rule means.Ruler, df map[string]any, fn func(Results, means.Ruler) []map[string]any) *Export {
	var keys []string
	for k := range df {
		keys = append(keys, k)
	}

	return &Export{
		rule:           rule,
		keys:           keys,
		defaultValues:  df,
		resultsToToken: fn,
	}
}

func NewExportDefault(rule means.Ruler, fn func(Results, means.Ruler) []map[string]any) *Export {
	return NewExport(rule, map[string]any{rule.NameAlias(): ""}, fn)
}

func (e *Export) GetKeys() []string {
	return e.keys
}

func (e *Export) GetDefaultValues() map[string]any {
	return e.defaultValues
}

func (e *Export) GetRule() means.Ruler {
	return e.rule
}

func (e *Export) Pluck(keys []string) *Export {
	df := maps.Clone(e.defaultValues)

	e.keys = make([]string, 0)
	e.defaultValues = make(map[string]any)

	for _, key := range keys {
		if v, ok := df[key]; ok {
			e.keys = append(e.keys, key)
			e.defaultValues[key] = v
		}
	}

	df = nil

	return e
}

func (e *Export) ToToken(results Results) search.Token {
	token := search.Token{}
	if results == nil {
		return token
	}

	token.SetOk()
	token.SetTokenizerFunc(func() []map[string]any {
		ret := e.resultsToToken(results, e.rule)

		// for pluck
		return maps2.PluckSliceMap(ret, e.keys)
	})

	return token
}

func NewExportAll(rule means.Ruler) *Export {
	df := map[string]any{
		rule.NameAlias():              "",
		rule.KeywordAmountNameAlias(): 0,
	}

	return NewExport(rule, df, func(results Results, rule means.Ruler) []map[string]any {
		if results == nil {
			return nil
		}

		return results.ToAll(rule)
	})
}

func NewExportLine(rule means.Ruler) *Export {
	df := map[string]any{
		rule.NameAlias():              "",
		rule.KeywordNumNameAlias():    0,
		rule.KeywordAmountNameAlias(): 0,
	}

	return NewExport(rule, df, func(results Results, rule means.Ruler) []map[string]any {
		if results == nil {
			return nil
		}

		return results.ToLine(rule)
	})
}

func NewExportFlag(rule means.Ruler) *Export {
	df := map[string]any{
		rule.NameAlias():        0,
		rule.KeywordNameAlias(): "",
	}

	return NewExport(rule, df, func(results Results, rule means.Ruler) []map[string]any {
		if results == nil {
			return nil
		}

		return results.ToFlag(rule)
	})
}
