package regexps

import (
	"maps"
	"strings"

	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
	maps2 "github.com/auho/go-etl/v2/tool/maps"
)

var _ search.Exporter = (*Export)(nil)

type Export struct {
	keys          []string
	defaultValues map[string]any

	resultsToToken func(Results, means.Ruler) []map[string]any
}

func NewExport(df map[string]any, fn func(Results, means.Ruler) []map[string]any) *Export {
	var keys []string
	for k := range df {
		keys = append(keys, k)
	}

	return &Export{
		keys:           keys,
		defaultValues:  df,
		resultsToToken: fn,
	}
}

func NewExportDefault(name string, fn func(Results, means.Ruler) []map[string]any) *Export {
	return NewExport(map[string]any{name: ""}, fn)
}

func (e *Export) GetKeys() []string {
	return e.keys
}

func (e *Export) GetDefaultValues() map[string]any {
	return e.defaultValues
}

func (e *Export) Pluck(keys []string) *Export {
	df := maps.Clone(e.defaultValues)

	e.keys = keys
	e.defaultValues = make(map[string]any, len(e.keys))

	for _, key := range e.keys {
		e.defaultValues[key] = df[key]
	}

	df = nil

	return e
}

func (e *Export) ToToken(results Results, rule means.Ruler) search.Token {
	token := search.Token{}
	if results == nil {
		return token
	}

	token.SetOk()
	token.SetTokenizerFunc(func() []map[string]any {
		ret := e.resultsToToken(results, rule)

		// for pluck
		return maps2.PluckSliceMap(ret, e.keys)
	})

	return token
}

func NewExportAll(rule means.Ruler) *Export {
	return NewExportDefault(rule.NameAlias(), func(results Results, rule means.Ruler) []map[string]any {
		if results == nil {
			return nil
		}

		var rets []map[string]any
		for _, result := range results {
			rets = append(rets, map[string]any{
				rule.NameAlias():              result.Text,
				rule.KeywordAmountNameAlias(): result.Amount,
			})
		}

		return rets
	})
}

func NewExportLine(rule means.Ruler) *Export {
	return NewExportDefault(rule.NameAlias(), func(results Results, rule means.Ruler) []map[string]any {
		if results == nil {
			return nil
		}

		var amount int
		var ss []string
		for _, result := range results {
			ss = append(ss, result.Text)
			amount += result.Amount
		}

		return []map[string]any{{
			rule.NameAlias():              strings.Join(ss, "|"),
			rule.KeywordAmountNameAlias(): amount,
		}}
	})
}

func NewExportFlag(rule means.Ruler) *Export {
	return NewExportDefault(rule.NameAlias(), func(results Results, rule means.Ruler) []map[string]any {
		if results == nil {
			return nil
		}

		var ss []string
		for _, result := range results {
			ss = append(ss, result.Text)
		}

		return []map[string]any{{rule.NameAlias(): 1, rule.KeywordNameAlias(): strings.Join(ss, "|")}}
	})
}
