package regexps

import (
	"strings"

	"github.com/auho/go-etl/v2/job/means"
)

type Result struct {
	Text   string
	Amount int
}

type Results []Result

func (rs Results) ToAll(rule means.Ruler) []map[string]any {
	var results []map[string]any
	for _, _r := range rs {
		results = append(results, map[string]any{
			rule.NameAlias():              _r.Text,
			rule.KeywordAmountNameAlias(): _r.Amount,
		})
	}

	return results
}

func (rs Results) ToLine(rule means.Ruler) []map[string]any {
	var ss []string
	var num, amount int
	for _, _r := range rs {
		ss = append(ss, _r.Text)
		num += 1
		amount += _r.Amount
	}

	return []map[string]any{
		{
			rule.NameAlias():              strings.Join(ss, "|"),
			rule.KeywordNumNameAlias():    num,
			rule.KeywordAmountNameAlias(): amount,
		},
	}
}

func (rs Results) ToFlag(rule means.Ruler) []map[string]any {
	var ss []string
	var num, amount int
	for _, _r := range rs {
		ss = append(ss, _r.Text)
		num += 1
		amount += _r.Amount
	}

	return []map[string]any{
		{
			rule.NameAlias():        1,
			rule.KeywordNameAlias(): strings.Join(ss, "|"),
		},
	}
}
