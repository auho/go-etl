package means

type TagMeans struct {
	TagMatcher
}

func (tm *TagMeans) insertResult(f func() *Result) [][]interface{} {
	result := f()
	if result == nil {
		return nil
	}

	return tm.ResultToSliceSlice(result)
}

func (tm *TagMeans) insertResults(f func() []*Result) [][]interface{} {
	results := f()
	if results == nil {
		return nil
	}

	return tm.ResultsToSliceSlice(results)
}

func (tm *TagMeans) updateResult(f func() *Result) map[string]interface{} {
	result := f()
	if result == nil {
		return nil
	}

	return tm.ResultToMap(result)
}

type TagMostTextMeans struct {
	TagMeans
}

func (t *TagMostTextMeans) GetKeys() []string {
	return t.GetResultInsertKeys()
}

func (t *TagMostTextMeans) Insert(contents []string) [][]interface{} {
	return t.insertResult(func() *Result {
		return t.Matcher.MatchMostText(contents)
	})
}

func (t *TagMostTextMeans) Update(contents []string) map[string]interface{} {
	return t.updateResult(func() *Result {
		return t.Matcher.MatchMostText(contents)
	})
}

type TagMostKeyMeans struct {
	TagMeans
}

func (t *TagMostKeyMeans) GetKeys() []string {
	return t.GetResultInsertKeys()
}

func (t *TagMostKeyMeans) Insert(contents []string) [][]interface{} {
	return t.insertResult(func() *Result {
		return t.Matcher.MatchMostKey(contents)
	})
}

func (t *TagMostKeyMeans) Update(contents []string) map[string]interface{} {
	return t.updateResult(func() *Result {
		return t.Matcher.MatchMostKey(contents)
	})
}
