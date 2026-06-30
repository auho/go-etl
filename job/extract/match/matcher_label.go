package match

import "github.com/auho/go-etl/v3/job/extract"

// label toMaps functions

func labelToMapsAll(r labelResults, rule extract.Rule, f Format) []map[string]any {
	return r.toAll(rule, f)
}

func labelToMapsLine(r labelResults, rule extract.Rule, f Format) []map[string]any {
	return r.toLine(rule, f)
}

func labelToMapsFlag(r labelResults, rule extract.Rule, f Format) []map[string]any {
	return r.toFlag(rule, f)
}

// label keysFun functions

func labelKeysAll(rule extract.Rule) ([]string, map[string]any) {
	keys := append(rule.TagsAlias(), rule.KeywordNameAlias(), rule.KeywordAmountNameAlias())
	defaults := make(map[string]any)
	for _, ta := range rule.TagsAlias() {
		defaults[ta] = ""
	}
	defaults[rule.KeywordNameAlias()] = ""
	defaults[rule.KeywordAmountNameAlias()] = 0
	return keys, defaults
}

func labelKeysLine(rule extract.Rule) ([]string, map[string]any) {
	keys := append(rule.TagsAlias(), rule.KeywordNameAlias(), rule.LabelNumNameAlias(), rule.KeywordNumNameAlias(), rule.KeywordAmountNameAlias())
	defaults := make(map[string]any)
	for _, ta := range rule.TagsAlias() {
		defaults[ta] = ""
	}
	defaults[rule.KeywordNameAlias()] = ""
	defaults[rule.LabelNumNameAlias()] = 0
	defaults[rule.KeywordNumNameAlias()] = 0
	defaults[rule.KeywordAmountNameAlias()] = 0
	return keys, defaults
}

func labelKeysFlag(rule extract.Rule) ([]string, map[string]any) {
	keys := append(rule.TagsAlias(), rule.KeywordNameAlias(), rule.NameAlias())
	defaults := make(map[string]any)
	for _, ta := range rule.TagsAlias() {
		defaults[ta] = ""
	}
	defaults[rule.KeywordNameAlias()] = ""
	defaults[rule.NameAlias()] = 0
	return keys, defaults
}

// constructors

func newMatcherLabel(rule extract.Rule, toMaps func(labelResults, extract.Rule, Format) []map[string]any, keysFun func(extract.Rule) ([]string, map[string]any), srf matcherResultsFunc[labelResults]) *MatcherLabelResults {
	return newMatcher[labelResults](rule, toMaps, keysFun, srf)
}

func newMatcherWholeLabels(rule extract.Rule, toMaps func(labelResults, extract.Rule, Format) []map[string]any, keysFun func(extract.Rule) ([]string, map[string]any)) *MatcherLabelResults {
	return newMatcherLabel(rule, toMaps, keysFun, func(ctx *matcherContextLabelResults, c []string) labelResults {
		return ctx.scanner.ScanLabel(c)
	})
}

func newMatcherLabels(rule extract.Rule, toMaps func(labelResults, extract.Rule, Format) []map[string]any, keysFun func(extract.Rule) ([]string, map[string]any)) *MatcherLabelResults {
	return newMatcherLabel(rule, toMaps, keysFun, func(ctx *matcherContextLabelResults, c []string) labelResults {
		return ctx.scanner.ScanLabel(c)
	})
}
