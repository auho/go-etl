package tag

import "github.com/auho/go-etl/v3/job/extract"

// label export configurations (toMaps / keys / defaults)
// mirroring the former NewExportLabel{All,Line,Flag} helpers.

func labelAllToMaps(r labelResults, rule extract.Rule, f format) []map[string]any {
	return r.toAll(rule, f)
}

func labelLineToMaps(r labelResults, rule extract.Rule, f format) []map[string]any {
	return r.toLine(rule, f)
}

func labelFlagToMaps(r labelResults, rule extract.Rule, f format) []map[string]any {
	return r.toFlag(rule, f)
}

func labelAllKeys(rule extract.Rule) []string {
	keys := make([]string, 0, len(rule.TagsAlias())+2)
	keys = append(keys, rule.TagsAlias()...)
	keys = append(keys, rule.KeywordNameAlias(), rule.KeywordAmountNameAlias())
	return keys
}

func labelLineKeys(rule extract.Rule) []string {
	keys := make([]string, 0, len(rule.TagsAlias())+4)
	keys = append(keys, rule.TagsAlias()...)
	keys = append(keys, rule.KeywordNameAlias(), rule.LabelNumNameAlias(), rule.KeywordNumNameAlias(), rule.KeywordAmountNameAlias())
	return keys
}

func labelFlagKeys(rule extract.Rule) []string {
	keys := make([]string, 0, len(rule.TagsAlias())+2)
	keys = append(keys, rule.TagsAlias()...)
	keys = append(keys, rule.KeywordNameAlias(), rule.NameAlias())
	return keys
}

func labelAllDefaults(rule extract.Rule) map[string]any {
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		values[_ta] = ""
	}
	values[rule.KeywordNameAlias()] = ""
	values[rule.KeywordAmountNameAlias()] = 0
	return values
}

func labelLineDefaults(rule extract.Rule) map[string]any {
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		values[_ta] = ""
	}
	values[rule.KeywordNameAlias()] = ""
	values[rule.LabelNumNameAlias()] = 0
	values[rule.KeywordNumNameAlias()] = 0
	values[rule.KeywordAmountNameAlias()] = 0
	return values
}

func labelFlagDefaults(rule extract.Rule) map[string]any {
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		values[_ta] = ""
	}
	values[rule.KeywordNameAlias()] = ""
	values[rule.NameAlias()] = 0
	return values
}

// constructors

func newSearchLabel(rule extract.Rule, toMaps func(labelResults, extract.Rule, format) []map[string]any, keys []string, defaults map[string]any, srf searchResultsFunc[labelResults]) *searchLabelResults {
	return newSearch[labelResults](rule, toMaps, keys, defaults, srf)
}

func newSearchWholeLabels(rule extract.Rule, toMaps func(labelResults, extract.Rule, format) []map[string]any, keys []string, defaults map[string]any) *searchLabelResults {
	return newSearchLabel(rule, toMaps, keys, defaults, func(ctx *searchContextLabelResults, c []string) labelResults {
		return ctx.matcher.MatchLabel(c)
	})
}

func newSearchLabels(rule extract.Rule, toMaps func(labelResults, extract.Rule, format) []map[string]any, keys []string, defaults map[string]any) *searchLabelResults {
	return newSearchLabel(rule, toMaps, keys, defaults, func(ctx *searchContextLabelResults, c []string) labelResults {
		return ctx.matcher.MatchLabel(c)
	})
}
