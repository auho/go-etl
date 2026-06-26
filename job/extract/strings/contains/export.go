package contains

import (
	"maps"

	"github.com/auho/go-etl/v3/job/extract"
	maps2 "github.com/auho/go-etl/v3/tool/mapx"
)

var _ extract.FieldSpec = (*Export)(nil)

type Export struct {
	rule           extract.Rule
	resultsToToken func(Results, extract.Rule) []map[string]any

	keys          []string
	defaultValues map[string]any
}

// NewExport
//
// df: map[string]any
// fn: func(Results, extract.Rule) []map[string]any
func NewExport(rule extract.Rule, df map[string]any, fn func(Results, extract.Rule) []map[string]any) *Export {
	var keys []string
	for k := range df {
		keys = append(keys, k)
	}

	return &Export{
		rule:           rule,
		resultsToToken: fn,
		keys:           keys,
		defaultValues:  df,
	}
}

func (e *Export) Keys() []string {
	return e.keys
}

func (e *Export) DefaultValues() map[string]any {
	return e.defaultValues
}

func (e *Export) GetRule() extract.Rule {
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

func (e *Export) ToToken(results Results) extract.Result {
	token := extract.Result{}

	if len(results) > 0 {
		token.SetOK()
		token.SetResultsFunc(func() []map[string]any {
			ret := e.resultsToToken(results, e.rule)

			// for pluck
			return maps2.PluckSliceMap(ret, e.keys)
		})
	}

	return token
}

func NewExportAll(rule extract.Rule) *Export {
	df := map[string]any{
		rule.NameAlias():              "",
		rule.KeywordAmountNameAlias(): 0,
	}

	return NewExport(rule, df, func(results Results, rule extract.Rule) []map[string]any {
		return results.ToAll(rule)
	})
}

func NewExportLine(rule extract.Rule) *Export {
	df := map[string]any{
		rule.NameAlias():              "",
		rule.KeywordNumNameAlias():    0,
		rule.KeywordAmountNameAlias(): 0,
	}

	return NewExport(rule, df, func(results Results, rule extract.Rule) []map[string]any {
		return results.ToLine(rule)
	})
}

func NewExportFlag(rule extract.Rule) *Export {
	df := map[string]any{
		rule.NameAlias():        0,
		rule.KeywordNameAlias(): "",
	}

	return NewExport(rule, df, func(results Results, rule extract.Rule) []map[string]any {
		return results.ToFlag(rule)
	})
}
