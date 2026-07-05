package tag

import "github.com/auho/go-etl/v3/task/extract"

// keyword export configurations (toMaps / keys / defaults)

func keywordAllToMaps(r results, rule extract.Rule, _ Format) []map[string]any {
	return r.toAll(rule)
}

func keywordLineToMaps(r results, rule extract.Rule, f Format) []map[string]any {
	return r.toLine(rule, f)
}

func keywordFlagToMaps(r results, rule extract.Rule, f Format) []map[string]any {
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

func newMatcherKeyword(rule extract.Rule, rowsFunc func(results, extract.Rule, Format) []map[string]any, keys []string, defaults map[string]any, rf resultsFunc[results]) *MatcherResults {
	return newMatcher[results](rule, rowsFunc, keys, defaults, rf)
}

func newMatcherFirstText(rule extract.Rule, rowsFunc func(results, extract.Rule, Format) []map[string]any, keys []string, defaults map[string]any) *MatcherResults {
	return newMatcherKeyword(rule, rowsFunc, keys, defaults, func(s *scanner, c []string) results {
		return s.ScanFirstText(c)
	})
}

func newMatcherLastText(rule extract.Rule, rowsFunc func(results, extract.Rule, Format) []map[string]any, keys []string, defaults map[string]any) *MatcherResults {
	return newMatcherKeyword(rule, rowsFunc, keys, defaults, func(s *scanner, c []string) results {
		return s.ScanLastText(c)
	})
}

func newMatcherMostText(rule extract.Rule, rowsFunc func(results, extract.Rule, Format) []map[string]any, keys []string, defaults map[string]any) *MatcherResults {
	return newMatcherKeyword(rule, rowsFunc, keys, defaults, func(s *scanner, c []string) results {
		return s.ScanMostText(c)
	})
}

func newMatcherKey(rule extract.Rule, rowsFunc func(results, extract.Rule, Format) []map[string]any, keys []string, defaults map[string]any) *MatcherResults {
	return newMatcherKeyword(rule, rowsFunc, keys, defaults, func(s *scanner, c []string) results {
		return s.ScanKey(c)
	})
}

func newMatcherFirstKey(rule extract.Rule, rowsFunc func(results, extract.Rule, Format) []map[string]any, keys []string, defaults map[string]any) *MatcherResults {
	return newMatcherKeyword(rule, rowsFunc, keys, defaults, func(s *scanner, c []string) results {
		return s.ScanFirstKey(c)
	})
}

func newMatcherLastKey(rule extract.Rule, rowsFunc func(results, extract.Rule, Format) []map[string]any, keys []string, defaults map[string]any) *MatcherResults {
	return newMatcherKeyword(rule, rowsFunc, keys, defaults, func(s *scanner, c []string) results {
		return s.ScanLastKey(c)
	})
}

func newMatcherMostKey(rule extract.Rule, rowsFunc func(results, extract.Rule, Format) []map[string]any, keys []string, defaults map[string]any) *MatcherResults {
	return newMatcherKeyword(rule, rowsFunc, keys, defaults, func(s *scanner, c []string) results {
		return s.ScanMostKey(c)
	})
}
