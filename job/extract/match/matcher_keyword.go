package match

import "github.com/auho/go-etl/v3/job/extract"

// keyword toMaps functions

func keywordToMapsAll(r results, rule extract.Rule, f Format) []map[string]any {
	return r.toAll(rule)
}

func keywordToMapsLine(r results, rule extract.Rule, f Format) []map[string]any {
	return r.toLine(rule, f)
}

func keywordToMapsFlag(r results, rule extract.Rule, f Format) []map[string]any {
	return r.toFlag(rule, f)
}

// keyword keysFun functions

func keywordKeysAll(rule extract.Rule) ([]string, map[string]any) {
	keys := append(rule.TagsAlias(), rule.KeywordNameAlias(), rule.KeywordAmountNameAlias())
	defaults := make(map[string]any)
	for _, ta := range rule.TagsAlias() {
		defaults[ta] = ""
	}
	defaults[rule.KeywordNameAlias()] = ""
	defaults[rule.KeywordAmountNameAlias()] = 0
	return keys, defaults
}

func keywordKeysLine(rule extract.Rule) ([]string, map[string]any) {
	keys := append(rule.TagsAlias(), rule.KeywordNameAlias(), rule.KeywordNumNameAlias())
	defaults := make(map[string]any)
	for _, ta := range rule.TagsAlias() {
		defaults[ta] = ""
	}
	defaults[rule.KeywordNameAlias()] = ""
	defaults[rule.KeywordNumNameAlias()] = 0
	return keys, defaults
}

func keywordKeysFlag(rule extract.Rule) ([]string, map[string]any) {
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

func newMatcherKeyword(rule extract.Rule, rowsFunc func(results, extract.Rule, Format) []map[string]any, keysFunc func(extract.Rule) ([]string, map[string]any), rf resultsFunc[results]) *MatcherResults {
	return newMatcher[results](rule, rowsFunc, keysFunc, rf)
}

func newMatcherFirstText(rule extract.Rule, rowsFunc func(results, extract.Rule, Format) []map[string]any, keysFunc func(extract.Rule) ([]string, map[string]any)) *MatcherResults {
	return newMatcherKeyword(rule, rowsFunc, keysFunc, func(s *scanner, c []string) results {
		return s.ScanFirstText(c)
	})
}

func newMatcherLastText(rule extract.Rule, rowsFunc func(results, extract.Rule, Format) []map[string]any, keysFunc func(extract.Rule) ([]string, map[string]any)) *MatcherResults {
	return newMatcherKeyword(rule, rowsFunc, keysFunc, func(s *scanner, c []string) results {
		return s.ScanLastText(c)
	})
}

func newMatcherMostText(rule extract.Rule, rowsFunc func(results, extract.Rule, Format) []map[string]any, keysFunc func(extract.Rule) ([]string, map[string]any)) *MatcherResults {
	return newMatcherKeyword(rule, rowsFunc, keysFunc, func(s *scanner, c []string) results {
		return s.ScanMostText(c)
	})
}

func newMatcherKey(rule extract.Rule, rowsFunc func(results, extract.Rule, Format) []map[string]any, keysFunc func(extract.Rule) ([]string, map[string]any)) *MatcherResults {
	return newMatcherKeyword(rule, rowsFunc, keysFunc, func(s *scanner, c []string) results {
		return s.ScanKey(c)
	})
}

func newMatcherFirstKey(rule extract.Rule, rowsFunc func(results, extract.Rule, Format) []map[string]any, keysFunc func(extract.Rule) ([]string, map[string]any)) *MatcherResults {
	return newMatcherKeyword(rule, rowsFunc, keysFunc, func(s *scanner, c []string) results {
		return s.ScanFirstKey(c)
	})
}

func newMatcherLastKey(rule extract.Rule, rowsFunc func(results, extract.Rule, Format) []map[string]any, keysFunc func(extract.Rule) ([]string, map[string]any)) *MatcherResults {
	return newMatcherKeyword(rule, rowsFunc, keysFunc, func(s *scanner, c []string) results {
		return s.ScanLastKey(c)
	})
}

func newMatcherMostKey(rule extract.Rule, rowsFunc func(results, extract.Rule, Format) []map[string]any, keysFunc func(extract.Rule) ([]string, map[string]any)) *MatcherResults {
	return newMatcherKeyword(rule, rowsFunc, keysFunc, func(s *scanner, c []string) results {
		return s.ScanMostKey(c)
	})
}
