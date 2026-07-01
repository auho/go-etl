package tag

import "github.com/auho/go-etl/v3/job/extract"

// label export configurations (toMaps / keys / defaults)

func labelAllToMaps(r labelResults, rule extract.Rule, f Format) []map[string]any {
	return r.toAll(rule, f)
}

func labelLineToMaps(r labelResults, rule extract.Rule, f Format) []map[string]any {
	return r.toLine(rule, f)
}

func labelFlagToMaps(r labelResults, rule extract.Rule, f Format) []map[string]any {
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

func newMatcherLabel(rule extract.Rule, rowsFunc func(labelResults, extract.Rule, Format) []map[string]any, keys []string, defaults map[string]any, rf resultsFunc[labelResults]) *MatcherLabelResults {
	return newMatcher[labelResults](rule, rowsFunc, keys, defaults, rf)
}

func newMatcherWholeLabels(rule extract.Rule, rowsFunc func(labelResults, extract.Rule, Format) []map[string]any, keys []string, defaults map[string]any) *MatcherLabelResults {
	return newMatcherLabel(rule, rowsFunc, keys, defaults, func(s *scanner, c []string) labelResults {
		return s.ScanLabel(c)
	})
}

func newMatcherLabels(rule extract.Rule, rowsFunc func(labelResults, extract.Rule, Format) []map[string]any, keys []string, defaults map[string]any) *MatcherLabelResults {
	return newMatcherLabel(rule, rowsFunc, keys, defaults, func(s *scanner, c []string) labelResults {
		return s.ScanLabel(c)
	})
}
