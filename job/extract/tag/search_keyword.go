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

func newSearchKeyword(rule extract.Rule, toMaps func(results, extract.Rule, format) []map[string]any, keys []string, defaults map[string]any, ste searchResultsFunc[results]) *searchResults {
	return newSearch[results](rule, toMaps, keys, defaults, ste)
}

func newSearchFirstText(rule extract.Rule, toMaps func(results, extract.Rule, format) []map[string]any, keys []string, defaults map[string]any) *searchResults {
	return newSearchKeyword(rule, toMaps, keys, defaults, func(ctx *searchContextResults, c []string) results {
		return ctx.matcher.MatchFirstText(c)
	})
}

func newSearchLastText(rule extract.Rule, toMaps func(results, extract.Rule, format) []map[string]any, keys []string, defaults map[string]any) *searchResults {
	return newSearchKeyword(rule, toMaps, keys, defaults, func(ctx *searchContextResults, c []string) results {
		return ctx.matcher.MatchLastText(c)
	})
}

func newSearchMostText(rule extract.Rule, toMaps func(results, extract.Rule, format) []map[string]any, keys []string, defaults map[string]any) *searchResults {
	return newSearchKeyword(rule, toMaps, keys, defaults, func(ctx *searchContextResults, c []string) results {
		return ctx.matcher.MatchMostText(c)
	})
}

func newSearchKey(rule extract.Rule, toMaps func(results, extract.Rule, format) []map[string]any, keys []string, defaults map[string]any) *searchResults {
	return newSearchKeyword(rule, toMaps, keys, defaults, func(ctx *searchContextResults, c []string) results {
		return ctx.matcher.MatchKey(c)
	})
}

func newSearchFirstKey(rule extract.Rule, toMaps func(results, extract.Rule, format) []map[string]any, keys []string, defaults map[string]any) *searchResults {
	return newSearchKeyword(rule, toMaps, keys, defaults, func(ctx *searchContextResults, c []string) results {
		return ctx.matcher.MatchFirstKey(c)
	})
}

func newSearchLastKey(rule extract.Rule, toMaps func(results, extract.Rule, format) []map[string]any, keys []string, defaults map[string]any) *searchResults {
	return newSearchKeyword(rule, toMaps, keys, defaults, func(ctx *searchContextResults, c []string) results {
		return ctx.matcher.MatchLastKey(c)
	})
}

func newSearchMostKey(rule extract.Rule, toMaps func(results, extract.Rule, format) []map[string]any, keys []string, defaults map[string]any) *searchResults {
	return newSearchKeyword(rule, toMaps, keys, defaults, func(ctx *searchContextResults, c []string) results {
		return ctx.matcher.MatchMostKey(c)
	})
}
