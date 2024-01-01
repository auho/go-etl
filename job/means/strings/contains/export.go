package contains

import (
	"maps"

	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
	maps2 "github.com/auho/go-etl/v2/tool/maps"
)

type Export struct {
	resultsToToken func(Results, means.Ruler) []map[string]any

	keys          []string
	defaultValues map[string]any
}

// NewExport
//
// keys: []string
// df: map[string]any
// fn: func(Results, means.Ruler) []map[string]any
func NewExport(keys []string, df map[string]any, fn func(Results, means.Ruler) []map[string]any) *Export {
	return &Export{
		resultsToToken: fn,
		keys:           keys,
		defaultValues:  df,
	}
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

func (e *Export) ToToken(rule means.Ruler, results Results) search.Token {
	token := search.Token{}

	if len(results) > 0 {
		token.SetOk()
		token.SetTokenizerFunc(func() []map[string]any {
			ret := e.resultsToToken(results, rule)

			// for pluck
			return maps2.PluckSliceMap(ret, e.keys)
		})
	}

	return token
}

func NewExportAll(rule means.Ruler) *Export {
	keys := []string{rule.NameAlias(), rule.KeywordAmountNameAlias()}
	df := map[string]any{
		rule.NameAlias():              "",
		rule.KeywordAmountNameAlias(): 0,
	}

	return NewExport(keys, df, func(results Results, rule means.Ruler) []map[string]any {
		return results.ToAll(rule)
	})
}

func NewExportLine(rule means.Ruler) *Export {
	keys := []string{rule.NameAlias(), rule.KeywordAmountNameAlias()}
	df := map[string]any{
		rule.NameAlias():              "",
		rule.KeywordNumNameAlias():    0,
		rule.KeywordAmountNameAlias(): 0,
	}

	return NewExport(keys, df, func(results Results, rule means.Ruler) []map[string]any {
		return results.ToLine(rule)
	})
}

func NewExportFlag(rule means.Ruler) *Export {
	keys := []string{rule.NameAlias(), rule.KeywordAmountNameAlias()}
	df := map[string]any{
		rule.NameAlias():        0,
		rule.KeywordNameAlias(): "",
	}

	return NewExport(keys, df, func(results Results, rule means.Ruler) []map[string]any {
		return results.ToFlag(rule)
	})
}
