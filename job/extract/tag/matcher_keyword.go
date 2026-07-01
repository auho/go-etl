package tag

import "github.com/auho/go-etl/v3/job/extract"

// keyword export configurations (toMaps / keys / defaults)
// mirroring the former NewExportKeyword{All,Line,Flag} helpers.

func keywordAllToMaps(r results, rule extract.Rule, _ format) []map[string]any {
	return r.toAll(rule)
}

func keywordLineToMaps(r results, rule extract.Rule, f format) []map[string]any {
	return r.toLine(rule, f)
}

func keywordFlagToMaps(r results, rule extract.Rule, f format) []map[string]any {
	return r.toFlag(rule, f)
}

func keywordAllKeys(rule extract.Rule) []string {
	keys := make([]string, 0, len(rule.TagsAlias())+2)
	keys = append(keys, rule.TagsAlias()...)
	keys = append(keys, rule.KeywordNameAlias(), rule.KeywordAmountNameAlias())
	return keys
}

func keywordLineKeys(rule extract.Rule) []string {
	keys := make([]string, 0, len(rule.TagsAlias())+2)
	keys = append(keys, rule.TagsAlias()...)
	keys = append(keys, rule.KeywordNameAlias(), rule.KeywordNumNameAlias())
	return keys
}

func keywordFlagKeys(rule extract.Rule) []string {
	keys := make([]string, 0, len(rule.TagsAlias())+2)
	keys = append(keys, rule.TagsAlias()...)
	keys = append(keys, rule.KeywordNameAlias(), rule.NameAlias())
	return keys
}

func keywordAllDefaults(rule extract.Rule) map[string]any {
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		values[_ta] = ""
	}
	values[rule.KeywordNameAlias()] = ""
	values[rule.KeywordAmountNameAlias()] = 0
	return values
}

func keywordLineDefaults(rule extract.Rule) map[string]any {
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		values[_ta] = ""
	}
	values[rule.KeywordNameAlias()] = ""
	values[rule.KeywordNumNameAlias()] = 0
	return values
}

func keywordFlagDefaults(rule extract.Rule) map[string]any {
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		values[_ta] = ""
	}
	values[rule.KeywordNameAlias()] = ""
	values[rule.NameAlias()] = 0
	return values
}

// constructors

func newMatcherKeyword(rule extract.Rule, toMaps func(results, extract.Rule, format) []map[string]any, keys []string, defaults map[string]any, ste matcherResultsFunc[results]) *MatcherResults {
	return newMatcher[results](rule, toMaps, keys, defaults, ste)
}

func newMatcherFirstText(rule extract.Rule, toMaps func(results, extract.Rule, format) []map[string]any, keys []string, defaults map[string]any) *MatcherResults {
	return newMatcherKeyword(rule, toMaps, keys, defaults, func(ctx *matcherContextResults, c []string) results {
		return ctx.scanner.ScanFirstText(c)
	})
}

func newMatcherLastText(rule extract.Rule, toMaps func(results, extract.Rule, format) []map[string]any, keys []string, defaults map[string]any) *MatcherResults {
	return newMatcherKeyword(rule, toMaps, keys, defaults, func(ctx *matcherContextResults, c []string) results {
		return ctx.scanner.ScanLastText(c)
	})
}

func newMatcherMostText(rule extract.Rule, toMaps func(results, extract.Rule, format) []map[string]any, keys []string, defaults map[string]any) *MatcherResults {
	return newMatcherKeyword(rule, toMaps, keys, defaults, func(ctx *matcherContextResults, c []string) results {
		return ctx.scanner.ScanMostText(c)
	})
}

func newMatcherKey(rule extract.Rule, toMaps func(results, extract.Rule, format) []map[string]any, keys []string, defaults map[string]any) *MatcherResults {
	return newMatcherKeyword(rule, toMaps, keys, defaults, func(ctx *matcherContextResults, c []string) results {
		return ctx.scanner.ScanKey(c)
	})
}

func newMatcherFirstKey(rule extract.Rule, toMaps func(results, extract.Rule, format) []map[string]any, keys []string, defaults map[string]any) *MatcherResults {
	return newMatcherKeyword(rule, toMaps, keys, defaults, func(ctx *matcherContextResults, c []string) results {
		return ctx.scanner.ScanFirstKey(c)
	})
}

func newMatcherLastKey(rule extract.Rule, toMaps func(results, extract.Rule, format) []map[string]any, keys []string, defaults map[string]any) *MatcherResults {
	return newMatcherKeyword(rule, toMaps, keys, defaults, func(ctx *matcherContextResults, c []string) results {
		return ctx.scanner.ScanLastKey(c)
	})
}

func newMatcherMostKey(rule extract.Rule, toMaps func(results, extract.Rule, format) []map[string]any, keys []string, defaults map[string]any) *MatcherResults {
	return newMatcherKeyword(rule, toMaps, keys, defaults, func(ctx *matcherContextResults, c []string) results {
		return ctx.scanner.ScanMostKey(c)
	})
}
