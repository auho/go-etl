package regexps

import (
	"strings"

	"github.com/auho/go-etl/v3/job/extract"
)

type result struct {
	text   string
	amount int
}

type results []result

func (rs results) toAll(rule extract.Rule) []map[string]any {
	var rets []map[string]any
	for _, _r := range rs {
		rets = append(rets, map[string]any{
			rule.NameAlias():              _r.text,
			rule.KeywordAmountNameAlias(): _r.amount,
		})
	}

	return rets
}

func (rs results) toLine(rule extract.Rule) []map[string]any {
	var ss []string
	var num, amount int
	for _, _r := range rs {
		ss = append(ss, _r.text)
		num += 1
		amount += _r.amount
	}

	return []map[string]any{
		{
			rule.NameAlias():              strings.Join(ss, "|"),
			rule.KeywordNumNameAlias():    num,
			rule.KeywordAmountNameAlias(): amount,
		},
	}
}

func (rs results) toFlag(rule extract.Rule) []map[string]any {
	var ss []string
	var num, amount int
	for _, _r := range rs {
		ss = append(ss, _r.text)
		num += 1
		amount += _r.amount
	}

	return []map[string]any{
		{
			rule.NameAlias():        1,
			rule.KeywordNameAlias(): strings.Join(ss, "|"),
		},
	}
}
