package match

import (
	"slices"
	"testing"

	"github.com/auho/go-etl/v3/internal/testutil"
	"github.com/auho/go-etl/v3/task/extract"
)

// TestExport_KeywordAll tests the full pipeline via NewKey
// (toAll)
func TestExport_KeywordAll(t *testing.T) {
	rule := &ruleTest{}
	s := NewKey(rule)
	s.scanner = _scanner
	defer s.Close()

	if err := s.Prepare(); err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}

	// Title should not be empty
	if s.Title() == "" {
		t.Fatal("Title() should not be empty")
	}

	// Keys should contain KeywordNameAlias and KeywordAmountNameAlias
	keys := s.Keys()
	if !slices.Contains(keys, rule.KeywordNameAlias()) {
		t.Errorf("Keys() should contain %s, got %v", rule.KeywordNameAlias(), keys)
	}
	if !slices.Contains(keys, rule.KeywordAmountNameAlias()) {
		t.Errorf("Keys() should contain %s, got %v", rule.KeywordAmountNameAlias(), keys)
	}

	// DefaultValues should be non-empty and contain expected keys
	dv := s.DefaultValues()
	if len(dv) == 0 {
		t.Error("DefaultValues() should not be empty")
	}
	if _, ok := dv[rule.KeywordNameAlias()]; !ok {
		t.Errorf("DefaultValues() should contain %s", rule.KeywordNameAlias())
	}
	if v := dv[rule.KeywordAmountNameAlias()]; v != 0 {
		t.Errorf("DefaultValues()[%s] should be 0, got %v", rule.KeywordAmountNameAlias(), v)
	}
	testutil.AssertMapCloned(t, "match.Matcher", s.DefaultValues)

	// Extract should return OK result with rows
	ret := s.Extract(_corpus)
	if !ret.IsOK() {
		t.Fatal("result.IsOK() should be true")
	}

	_, rows := ret.Get()
	if len(rows) == 0 {
		t.Fatal("rows should not be empty")
	}

	// Each row should contain keyword and amount keys (ToAll keys)
	for i, row := range rows {
		if _, ok := row[rule.KeywordNameAlias()]; !ok {
			t.Errorf("row %d missing key %s", i, rule.KeywordNameAlias())
		}
		if _, ok := row[rule.KeywordAmountNameAlias()]; !ok {
			t.Errorf("row %d missing key %s", i, rule.KeywordAmountNameAlias())
		}
	}
}

// TestExport_KeywordLine tests keyword line export with newMatcherKey
// (toLine)
func TestExport_KeywordLine(t *testing.T) {
	rule := &ruleTest{}
	s := newMatcherKey(rule, keywordToMapsLine, keywordKeysLine)
	s.scanner = _scanner
	defer s.Close()

	if err := s.Prepare(); err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}

	ret := s.Extract(_corpus)
	if !ret.IsOK() {
		t.Fatal("result.IsOK() should be true")
	}

	// ToLine merges all results into a single row
	_, rows := ret.Get()
	if len(rows) != 1 {
		t.Fatalf("ToLine should return 1 row, got %d", len(rows))
	}

	// Line keys: TagsAlias + KeywordNameAlias + KeywordNumNameAlias
	if _, ok := rows[0][rule.KeywordNameAlias()]; !ok {
		t.Errorf("row missing key %s", rule.KeywordNameAlias())
	}
	if _, ok := rows[0][rule.KeywordNumNameAlias()]; !ok {
		t.Errorf("row missing key %s", rule.KeywordNumNameAlias())
	}
}

// TestExport_KeywordFlag tests keyword flag export with newMatcherKey
// (toFlag)
func TestExport_KeywordFlag(t *testing.T) {
	rule := &ruleTest{}
	s := newMatcherKey(rule, keywordToMapsFlag, keywordKeysFlag)
	s.scanner = _scanner
	defer s.Close()

	if err := s.Prepare(); err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}

	ret := s.Extract(_corpus)
	if !ret.IsOK() {
		t.Fatal("result.IsOK() should be true")
	}

	// ToFlag merges all results into a single row with a flag
	_, rows := ret.Get()
	if len(rows) != 1 {
		t.Fatalf("ToFlag should return 1 row, got %d", len(rows))
	}

	// Flag keys: TagsAlias + KeywordNameAlias
	if _, ok := rows[0][rule.KeywordNameAlias()]; !ok {
		t.Errorf("row missing key %s", rule.KeywordNameAlias())
	}
}

// TestExport_LabelAll tests the full pipeline via NewLabel
// (labelResults.toAll)
func TestExport_LabelAll(t *testing.T) {
	rule := &ruleTest{}
	s := NewLabel(rule)
	s.scanner = _scanner
	defer s.Close()

	if err := s.Prepare(); err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}

	ret := s.Extract(_corpus)
	if !ret.IsOK() {
		t.Fatal("result.IsOK() should be true")
	}

	_, rows := ret.Get()
	if len(rows) == 0 {
		t.Fatal("rows should not be empty")
	}

	for i, row := range rows {
		if _, ok := row[rule.KeywordNameAlias()]; !ok {
			t.Errorf("row %d missing key %s", i, rule.KeywordNameAlias())
		}
		if _, ok := row[rule.KeywordAmountNameAlias()]; !ok {
			t.Errorf("row %d missing key %s", i, rule.KeywordAmountNameAlias())
		}
	}
}

// TestExport_LabelLine tests the full pipeline via NewWholeLabels
// (labelResults.toLine / mergeLabelsToWhole)
func TestExport_LabelLine(t *testing.T) {
	rule := &ruleTest{}
	s := NewWholeLabels(rule)
	s.scanner = _scanner
	defer s.Close()

	if err := s.Prepare(); err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}

	ret := s.Extract(_corpus)
	if !ret.IsOK() {
		t.Fatal("result.IsOK() should be true")
	}

	// ToLine merges all labels into a single row
	_, rows := ret.Get()
	if len(rows) != 1 {
		t.Fatalf("ToLine should return 1 row, got %d", len(rows))
	}

	// Line keys: TagsAlias + KeywordNameAlias + LabelNumNameAlias + KeywordNumNameAlias + KeywordAmountNameAlias
	if _, ok := rows[0][rule.KeywordNameAlias()]; !ok {
		t.Errorf("row missing key %s", rule.KeywordNameAlias())
	}
	if _, ok := rows[0][rule.LabelNumNameAlias()]; !ok {
		t.Errorf("row missing key %s", rule.LabelNumNameAlias())
	}
}

// TestExport_LabelFlag tests label flag export with newMatcherLabels
// (labelResults.toFlag / mergeLabelsToWhole)
func TestExport_LabelFlag(t *testing.T) {
	rule := &ruleTest{}
	s := newMatcherLabels(rule, labelToMapsFlag, labelKeysFlag)
	s.scanner = _scanner
	defer s.Close()

	if err := s.Prepare(); err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}

	ret := s.Extract(_corpus)
	if !ret.IsOK() {
		t.Fatal("result.IsOK() should be true")
	}

	_, rows := ret.Get()
	if len(rows) != 1 {
		t.Fatalf("ToFlag should return 1 row, got %d", len(rows))
	}

	// Flag keys: TagsAlias + KeywordNameAlias
	if _, ok := rows[0][rule.KeywordNameAlias()]; !ok {
		t.Errorf("row missing key %s", rule.KeywordNameAlias())
	}
}

// TestExport_Pluck tests WithPluck functionality
func TestExport_Pluck(t *testing.T) {
	rule := &ruleTest{}
	pluckedKey := rule.KeywordNameAlias()
	s := NewKey(rule).WithPluckKeys([]string{pluckedKey})
	s.scanner = _scanner
	defer s.Close()

	if err := s.Prepare(); err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}

	// After pluck, Keys should only contain the plucked key
	keys := s.Keys()
	if len(keys) != 1 {
		t.Fatalf("expected 1 key after pluck, got %d: %v", len(keys), keys)
	}
	if keys[0] != pluckedKey {
		t.Errorf("expected key %s, got %s", pluckedKey, keys[0])
	}

	// DefaultValues should only contain the plucked key
	dv := s.DefaultValues()
	if len(dv) != 1 {
		t.Fatalf("expected 1 default value, got %d", len(dv))
	}
	if _, ok := dv[pluckedKey]; !ok {
		t.Errorf("DefaultValues should contain %s", pluckedKey)
	}

	ret := s.Extract(_corpus)
	if !ret.IsOK() {
		t.Fatal("result.IsOK() should be true")
	}

	_, rows := ret.Get()
	if len(rows) == 0 {
		t.Fatal("rows should not be empty")
	}

	// Each row should only have the plucked key
	for i, row := range rows {
		if len(row) != 1 {
			t.Errorf("row %d should have 1 key, got %d: %v", i, len(row), row)
		}
		if _, ok := row[pluckedKey]; !ok {
			t.Errorf("row %d should contain %s", i, pluckedKey)
		}
	}
}

// TestExport_WithFormat tests WithFormat on Matcher
func TestExport_WithFormat(t *testing.T) {
	rule := &ruleTest{}
	customFormat := Format{withKeywordAmount: false, sep: "|"}
	s := newMatcherKey(rule, keywordToMapsLine, keywordKeysLine).
		WithFormat(customFormat)
	s.scanner = _scanner
	defer s.Close()

	if err := s.Prepare(); err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}

	ret := s.Extract(_corpus)
	if !ret.IsOK() {
		t.Fatal("result.IsOK() should be true")
	}

	_, rows := ret.Get()
	if len(rows) != 1 {
		t.Fatalf("ToLine should return 1 row, got %d", len(rows))
	}

	// WithKeywordAmount=false: keyword text has no amount suffix
	// Sep="|": keywords are joined by |
	keywordValue, ok := rows[0][rule.KeywordNameAlias()].(string)
	if !ok {
		t.Fatalf("expected string for %s, got %T", rule.KeywordNameAlias(), rows[0][rule.KeywordNameAlias()])
	}

	if keywordValue == "" {
		t.Error("keyword value should not be empty")
	}
}

// TestExport_EmptyResults tests no-match scenario
// Uses ruleTest's defaultScanner (keywords: "123", "b", "e", "中文", "中1文")
// with content that contains none of these keywords
func TestExport_EmptyResults(t *testing.T) {
	rule := &ruleTest{}
	s := NewKey(rule)
	defer s.Close()

	if err := s.Prepare(); err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}

	ret := s.Extract([]string{"zzzzzzzzzzz9999999999"})
	if ret.IsOK() {
		t.Fatal("result.IsOK() should be false for no match")
	}

	_, rows := ret.Get()
	if rows != nil {
		t.Fatalf("rows should be nil for no match, got %v", rows)
	}
}

// TestExport_AllEntryFunctions tests all entry.go entry functions
func TestExport_AllEntryFunctions(t *testing.T) {
	rule := &ruleTest{}

	entryFuncs := []struct {
		name string
		fn   func(extract.Rule) *MatcherResults
	}{
		{"NewKey", NewKey},
		{"NewMostKey", NewMostKey},
		{"NewMostText", NewMostText},
		{"NewFirstText", NewFirstText},
		{"NewFirstKey", NewFirstKey},
	}

	for _, ef := range entryFuncs {
		t.Run(ef.name, func(t *testing.T) {
			s := ef.fn(rule)
			s.scanner = _scanner
			defer s.Close()

			if err := s.Prepare(); err != nil {
				t.Fatalf("Prepare failed: %v", err)
			}

			ret := s.Extract(_corpus)
			_ = ret.IsOK()
			_, _ = ret.Get()
		})
	}

	labelFuncs := []struct {
		name string
		fn   func(extract.Rule) *MatcherLabelResults
	}{
		{"NewLabel", NewLabel},
		{"NewWholeLabels", NewWholeLabels},
	}

	for _, lf := range labelFuncs {
		t.Run(lf.name, func(t *testing.T) {
			s := lf.fn(rule)
			s.scanner = _scanner
			defer s.Close()

			if err := s.Prepare(); err != nil {
				t.Fatalf("Prepare failed: %v", err)
			}

			ret := s.Extract(_corpus)
			_ = ret.IsOK()
			_, _ = ret.Get()
		})
	}
}

// TestExport_ExtraMatcherFunctions tests newMatcherLastText and newMatcherLastKey
// which are not covered by entry.go entry functions
func TestExport_ExtraMatcherFunctions(t *testing.T) {
	rule := &ruleTest{}

	t.Run("newMatcherLastText", func(t *testing.T) {
		s := newMatcherLastText(rule, keywordToMapsAll, keywordKeysAll)
		s.scanner = _scanner
		defer s.Close()

		if err := s.Prepare(); err != nil {
			t.Fatalf("Prepare failed: %v", err)
		}

		ret := s.Extract(_corpus)
		_ = ret.IsOK()
		_, _ = ret.Get()
	})

	t.Run("newMatcherLastKey", func(t *testing.T) {
		s := newMatcherLastKey(rule, keywordToMapsAll, keywordKeysAll)
		s.scanner = _scanner
		defer s.Close()

		if err := s.Prepare(); err != nil {
			t.Fatalf("Prepare failed: %v", err)
		}

		ret := s.Extract(_corpus)
		_ = ret.IsOK()
		_, _ = ret.Get()
	})
}

// TestExport_MatchOptions tests match-specific Matcher options:
// WithIgnoreCase, WithFuzzy, WithPriorityFuzzy, WithDebug
func TestExport_MatchOptions(t *testing.T) {
	rule := &ruleTest{}

	t.Run("WithIgnoreCase", func(t *testing.T) {
		s := NewKey(rule).
			WithIgnoreCase()
		s.scanner = _scanner
		defer s.Close()

		if err := s.Prepare(); err != nil {
			t.Fatalf("Prepare failed: %v", err)
		}

		ret := s.Extract(_corpus)
		_ = ret.IsOK()
		_, _ = ret.Get()
	})

	t.Run("WithFuzzy", func(t *testing.T) {
		s := NewKey(rule).
			WithFuzzy(FuzzyConfig{Window: 3, Sep: "_"})
		s.scanner = _scanner
		defer s.Close()

		if err := s.Prepare(); err != nil {
			t.Fatalf("Prepare failed: %v", err)
		}

		ret := s.Extract(_corpus)
		_ = ret.IsOK()
		_, _ = ret.Get()
	})

	t.Run("WithPriorityFuzzy", func(t *testing.T) {
		s := NewKey(rule).
			WithModePriorityFuzzy()
		defer s.Close()

		if err := s.Prepare(); err != nil {
			t.Fatalf("Prepare failed: %v", err)
		}

		ret := s.Extract(_corpus)
		_ = ret.IsOK()
		_, _ = ret.Get()
	})

	t.Run("WithDebug", func(t *testing.T) {
		s := NewKey(rule).
			WithDebug()
		s.scanner = _scanner
		defer s.Close()

		if err := s.Prepare(); err != nil {
			t.Fatalf("Prepare failed: %v", err)
		}

		ret := s.Extract(_corpus)
		_ = ret.IsOK()
		_, _ = ret.Get()
	})

	t.Run("CombinedOptions", func(t *testing.T) {
		s := NewKey(rule).
			WithIgnoreCase().
			WithFuzzy(FuzzyConfig{Window: 3, Sep: "_"}).
			WithModePriorityFuzzy().
			WithDebug()
		s.scanner = _scanner
		defer s.Close()

		if err := s.Prepare(); err != nil {
			t.Fatalf("Prepare failed: %v", err)
		}

		ret := s.Extract(_corpus)
		_ = ret.IsOK()
		_, _ = ret.Get()
	})
}

// TestExport_Constructors tests newMatcher with custom rowsFunc/keysFunc and WithPluck
func TestExport_Constructors(t *testing.T) {
	rule := &ruleTest{}

	t.Run("newMatcher_Keyword", func(t *testing.T) {
		keysFunc := func(r extract.Rule) ([]string, map[string]any) {
			return []string{"k1", "k2"}, map[string]any{"k1": "v1", "k2": 0}
		}
		rowsFunc := func(r results, rule extract.Rule, f Format) []map[string]any {
			return []map[string]any{{"k1": "x", "k2": 1}}
		}
		s := newMatcher[results](rule, rowsFunc, keysFunc, func(s *scanner, c []string) results {
			return nil
		})

		if err := s.Prepare(); err != nil {
			t.Fatalf("Prepare failed: %v", err)
		}

		if len(s.Keys()) != 2 {
			t.Errorf("expected 2 keys, got %d", len(s.Keys()))
		}

		if len(s.DefaultValues()) != 2 {
			t.Errorf("expected 2 default values, got %d", len(s.DefaultValues()))
		}
	})

	t.Run("newMatcher_Label", func(t *testing.T) {
		keysFunc := func(r extract.Rule) ([]string, map[string]any) {
			return []string{"k1"}, map[string]any{"k1": "v1"}
		}
		rowsFunc := func(r labelResults, rule extract.Rule, f Format) []map[string]any {
			return []map[string]any{{"k1": "x"}}
		}
		s := newMatcher[labelResults](rule, rowsFunc, keysFunc, func(s *scanner, c []string) labelResults {
			return nil
		})

		if err := s.Prepare(); err != nil {
			t.Fatalf("Prepare failed: %v", err)
		}

		if len(s.Keys()) != 1 {
			t.Errorf("expected 1 key, got %d", len(s.Keys()))
		}
	})

	t.Run("Pluck_Direct", func(t *testing.T) {
		keysFunc := func(r extract.Rule) ([]string, map[string]any) {
			return []string{"k1", "k2", "k3"}, map[string]any{"k1": "v1", "k2": 0, "k3": "v3"}
		}
		rowsFunc := func(r results, rule extract.Rule, f Format) []map[string]any {
			return nil
		}
		s := newMatcher[results](rule, rowsFunc, keysFunc, func(s *scanner, c []string) results {
			return nil
		})
		s.WithPluckKeys([]string{"k1", "k3"})

		if err := s.Prepare(); err != nil {
			t.Fatalf("Prepare failed: %v", err)
		}

		keys := s.Keys()
		if len(keys) != 2 {
			t.Fatalf("expected 2 keys after pluck, got %d", len(keys))
		}
		if !slices.Contains(keys, "k1") || !slices.Contains(keys, "k3") {
			t.Errorf("expected keys k1 and k3, got %v", keys)
		}

		dv := s.DefaultValues()
		if len(dv) != 2 {
			t.Fatalf("expected 2 default values after pluck, got %d", len(dv))
		}
		if dv["k1"] != "v1" {
			t.Errorf("expected v1 for k1, got %v", dv["k1"])
		}
	})

	t.Run("Pluck_NonExistentKey", func(t *testing.T) {
		keysFunc := func(r extract.Rule) ([]string, map[string]any) {
			return []string{"k1"}, map[string]any{"k1": "v1"}
		}
		rowsFunc := func(r results, rule extract.Rule, f Format) []map[string]any {
			return nil
		}
		s := newMatcher[results](rule, rowsFunc, keysFunc, func(s *scanner, c []string) results {
			return nil
		})
		s.WithPluckKeys([]string{"k1", "nonexistent"})

		if err := s.Prepare(); err != nil {
			t.Fatalf("Prepare failed: %v", err)
		}

		keys := s.Keys()
		if len(keys) != 1 {
			t.Fatalf("expected 1 key after pluck, got %d", len(keys))
		}
		if keys[0] != "k1" {
			t.Errorf("expected key k1, got %s", keys[0])
		}
	})
}

// TestExport_ResultConstructors tests newResult, newLabelResult, toTag methods
func TestExport_ResultConstructors(t *testing.T) {
	rule := &ruleTest{}

	t.Run("NewResult", func(t *testing.T) {
		r := newResult()
		if r.tags == nil {
			t.Error("Tags should not be nil")
		}
		if r.texts == nil {
			t.Error("Texts should not be nil")
		}
	})

	t.Run("newLabelResult", func(t *testing.T) {
		lr := newLabelResult()
		if lr.tags == nil {
			t.Error("Tags should not be nil")
		}
		if lr.match == nil {
			t.Error("Match should not be nil")
		}
	})

	t.Run("Result_ToTag", func(t *testing.T) {
		r := newResult()
		r.keyword = "test_keyword"
		r.amount = 5
		r.tags["a"] = "tag_a"
		r.tags["ab"] = "tag_ab"
		r.texts["test_keyword"] = 5

		tag := r.toTag(rule)
		if tag[rule.KeywordNameAlias()] != "test_keyword" {
			t.Errorf("expected keyword test_keyword, got %v", tag[rule.KeywordNameAlias()])
		}
		if tag[rule.KeywordAmountNameAlias()] != 5 {
			t.Errorf("expected amount 5, got %v", tag[rule.KeywordAmountNameAlias()])
		}
		if tag[rule.KeywordNumNameAlias()] != 1 {
			t.Errorf("expected num 1, got %v", tag[rule.KeywordNumNameAlias()])
		}
		if tag["a"] != "tag_a" {
			t.Errorf("expected tag_a, got %v", tag["a"])
		}
	})

	t.Run("LabelResult_ToTag", func(t *testing.T) {
		lr := newLabelResult()
		lr.identity = "test"
		lr.amount = 3
		lr.tags["a"] = "tag_a"
		lr.match["key1"] = map[string]int{"text1": 2, "text2": 1}
		lr.keywords = []string{"key1"}

		ltag := lr.toTag(rule, defaultFormat)
		// WithKeywordAmount=true, keyText = "key1 3" (key + total amount)
		if ltag[rule.KeywordNameAlias()] != "key1 3" {
			t.Errorf("expected 'key1 3', got %v", ltag[rule.KeywordNameAlias()])
		}
		if ltag[rule.KeywordNumNameAlias()] != 1 {
			t.Errorf("expected num 1, got %v", ltag[rule.KeywordNumNameAlias()])
		}
		if ltag[rule.KeywordAmountNameAlias()] != 3 {
			t.Errorf("expected amount 3, got %v", ltag[rule.KeywordAmountNameAlias()])
		}
	})

	t.Run("LabelResult_ToTag_NoAmount", func(t *testing.T) {
		lr := newLabelResult()
		lr.amount = 3
		lr.match["key1"] = map[string]int{"text1": 2, "text2": 1}
		lr.keywords = []string{"key1"}

		format := Format{withKeywordAmount: false, sep: ","}
		ltag := lr.toTag(rule, format)
		// WithKeywordAmount=false, keyText = "key1" (no amount)
		if ltag[rule.KeywordNameAlias()] != "key1" {
			t.Errorf("expected 'key1', got %v", ltag[rule.KeywordNameAlias()])
		}
	})
}
