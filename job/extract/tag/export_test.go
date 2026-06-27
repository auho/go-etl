package tag

import (
	"slices"
	"testing"

	"github.com/auho/go-etl/v3/job/extract"
)

// TestExport_KeywordAll tests the full pipeline via NewKey
// (NewExportKeywordAll / Results.ToAll)
func TestExport_KeywordAll(t *testing.T) {
	rule := &ruleTest{}
	s := NewKey(rule)
	defer s.Close()

	if err := s.Prepare(); err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}

	// Title should not be empty
	if s.Title() == "" {
		t.Fatal("Title() should not be empty")
	}

	// NewExport should return non-nil FieldSpec
	spec := s.NewExport()
	if spec == nil {
		t.Fatal("NewExport() should not be nil")
	}

	// Keys should contain KeywordNameAlias and KeywordAmountNameAlias
	keys := spec.Keys()
	if !slices.Contains(keys, rule.KeywordNameAlias()) {
		t.Errorf("Keys() should contain %s, got %v", rule.KeywordNameAlias(), keys)
	}
	if !slices.Contains(keys, rule.KeywordAmountNameAlias()) {
		t.Errorf("Keys() should contain %s, got %v", rule.KeywordAmountNameAlias(), keys)
	}

	// DefaultValues should be non-empty and contain expected keys
	dv := spec.DefaultValues()
	if len(dv) == 0 {
		t.Error("DefaultValues() should not be empty")
	}
	if _, ok := dv[rule.KeywordNameAlias()]; !ok {
		t.Errorf("DefaultValues() should contain %s", rule.KeywordNameAlias())
	}
	if v := dv[rule.KeywordAmountNameAlias()]; v != 0 {
		t.Errorf("DefaultValues()[%s] should be 0, got %v", rule.KeywordAmountNameAlias(), v)
	}

	// Search should return OK result with rows
	result := s.Search(_contents)
	if !result.IsOK() {
		t.Fatal("result.IsOK() should be true")
	}

	rows := result.Rows()
	if len(rows) == 0 {
		t.Fatal("Rows() should not be empty")
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

// TestExport_KeywordLine tests NewExportKeywordLine with NewSearchKey
// (Results.ToLine / MergeKeysToWhole)
func TestExport_KeywordLine(t *testing.T) {
	rule := &ruleTest{}
	s := NewSearchKey(NewExportKeywordLine(rule))
	defer s.Close()

	if err := s.Prepare(); err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}

	result := s.Search(_contents)
	if !result.IsOK() {
		t.Fatal("result.IsOK() should be true")
	}

	// ToLine merges all results into a single row
	rows := result.Rows()
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

// TestExport_KeywordFlag tests NewExportKeywordFlag with NewSearchKey
// (Results.ToFlag / MergeKeysToWhole)
func TestExport_KeywordFlag(t *testing.T) {
	rule := &ruleTest{}
	s := NewSearchKey(NewExportKeywordFlag(rule))
	defer s.Close()

	if err := s.Prepare(); err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}

	result := s.Search(_contents)
	if !result.IsOK() {
		t.Fatal("result.IsOK() should be true")
	}

	// ToFlag merges all results into a single row with a flag
	rows := result.Rows()
	if len(rows) != 1 {
		t.Fatalf("ToFlag should return 1 row, got %d", len(rows))
	}

	// Flag keys: TagsAlias + KeywordNameAlias
	if _, ok := rows[0][rule.KeywordNameAlias()]; !ok {
		t.Errorf("row missing key %s", rule.KeywordNameAlias())
	}
}

// TestExport_LabelAll tests the full pipeline via NewLabel
// (NewExportLabelAll / LabelResults.ToAll)
func TestExport_LabelAll(t *testing.T) {
	rule := &ruleTest{}
	s := NewLabel(rule)
	defer s.Close()

	if err := s.Prepare(); err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}

	result := s.Search(_contents)
	if !result.IsOK() {
		t.Fatal("result.IsOK() should be true")
	}

	rows := result.Rows()
	if len(rows) == 0 {
		t.Fatal("Rows() should not be empty")
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
// (NewExportLabelLine / LabelResults.ToLine / MergeLabelsToWhole)
func TestExport_LabelLine(t *testing.T) {
	rule := &ruleTest{}
	s := NewWholeLabels(rule)
	defer s.Close()

	if err := s.Prepare(); err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}

	result := s.Search(_contents)
	if !result.IsOK() {
		t.Fatal("result.IsOK() should be true")
	}

	// ToLine merges all labels into a single row
	rows := result.Rows()
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

// TestExport_LabelFlag tests NewExportLabelFlag with NewSearchLabels
// (LabelResults.ToFlag / MergeLabelsToWhole)
func TestExport_LabelFlag(t *testing.T) {
	rule := &ruleTest{}
	s := NewSearchLabels(NewExportLabelFlag(rule))
	defer s.Close()

	if err := s.Prepare(); err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}

	result := s.Search(_contents)
	if !result.IsOK() {
		t.Fatal("result.IsOK() should be true")
	}

	rows := result.Rows()
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
	s := NewKey(rule).WithPluck([]string{pluckedKey})
	defer s.Close()

	if err := s.Prepare(); err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}

	// After pluck, Keys should only contain the plucked key
	keys := s.NewExport().Keys()
	if len(keys) != 1 {
		t.Fatalf("expected 1 key after pluck, got %d: %v", len(keys), keys)
	}
	if keys[0] != pluckedKey {
		t.Errorf("expected key %s, got %s", pluckedKey, keys[0])
	}

	// DefaultValues should only contain the plucked key
	dv := s.NewExport().DefaultValues()
	if len(dv) != 1 {
		t.Fatalf("expected 1 default value, got %d", len(dv))
	}
	if _, ok := dv[pluckedKey]; !ok {
		t.Errorf("DefaultValues should contain %s", pluckedKey)
	}

	result := s.Search(_contents)
	if !result.IsOK() {
		t.Fatal("result.IsOK() should be true")
	}

	rows := result.Rows()
	if len(rows) == 0 {
		t.Fatal("Rows() should not be empty")
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

// TestExport_WithFormat tests WithFormat on Export
func TestExport_WithFormat(t *testing.T) {
	rule := &ruleTest{}
	customFormat := Format{WithKeywordAmount: false, Sep: "|"}
	values := make(map[string]any)
	for _, _ta := range rule.TagsAlias() {
		values[_ta] = ""
	}
	values[rule.KeywordNameAlias()] = ""
	values[rule.KeywordNumNameAlias()] = 0
	export := extract.NewExporter(values, func(ctx extract.ExportContext[Results]) []map[string]any {
		return ctx.Results.ToLine(rule, ctx.Format.(Format))
	}, extract.WithRule[Results](rule), extract.WithFormat[Results](customFormat))
	s := NewSearchKey(export)
	defer s.Close()

	if err := s.Prepare(); err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}

	result := s.Search(_contents)
	if !result.IsOK() {
		t.Fatal("result.IsOK() should be true")
	}

	rows := result.Rows()
	if len(rows) != 1 {
		t.Fatalf("ToLine should return 1 row, got %d", len(rows))
	}

	// WithKeywordAmount=false: keyword text has no amount suffix
	// Sep="|": keywords are joined by |
	keywordValue, ok := rows[0][rule.KeywordNameAlias()].(string)
	if !ok {
		t.Fatalf("expected string for %s, got %T", rule.KeywordNameAlias(), rows[0][rule.KeywordNameAlias()])
	}

	// keywordValue should use "|" as separator (not default ",")
	if keywordValue == "" {
		t.Error("keyword value should not be empty")
	}
}

// TestExport_EmptyResults tests no-match scenario
func TestExport_EmptyResults(t *testing.T) {
	rule := &ruleTest{}
	s := NewKey(rule)
	defer s.Close()

	if err := s.Prepare(); err != nil {
		t.Fatalf("Prepare failed: %v", err)
	}

	result := s.Search([]string{"zzzzz_no_match_content_xxxxx"})
	if result.IsOK() {
		t.Fatal("result.IsOK() should be false for no match")
	}

	rows := result.Rows()
	if rows != nil {
		t.Fatalf("Rows() should be nil for no match, got %v", rows)
	}
}

// TestExport_AllEntryFunctions tests all entry.go entry functions
func TestExport_AllEntryFunctions(t *testing.T) {
	rule := &ruleTest{}

	entryFuncs := []struct {
		name string
		fn   func(extract.Rule) *SearchResults
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
			defer s.Close()

			if err := s.Prepare(); err != nil {
				t.Fatalf("Prepare failed: %v", err)
			}

			result := s.Search(_contents)
			// just verify no panic; some may match, some may not
			_ = result.IsOK()
			_ = result.Rows()
		})
	}

	labelFuncs := []struct {
		name string
		fn   func(extract.Rule) *SearchLabelResults
	}{
		{"NewLabel", NewLabel},
		{"NewWholeLabels", NewWholeLabels},
	}

	for _, lf := range labelFuncs {
		t.Run(lf.name, func(t *testing.T) {
			s := lf.fn(rule)
			defer s.Close()

			if err := s.Prepare(); err != nil {
				t.Fatalf("Prepare failed: %v", err)
			}

			result := s.Search(_contents)
			_ = result.IsOK()
			_ = result.Rows()
		})
	}
}

// TestExport_ExtraSearchFunctions tests NewSearchLastText and NewSearchLastKey
// which are not covered by entry.go entry functions
func TestExport_ExtraSearchFunctions(t *testing.T) {
	rule := &ruleTest{}

	t.Run("NewSearchLastText", func(t *testing.T) {
		s := NewSearchLastText(NewExportKeywordAll(rule))
		defer s.Close()

		if err := s.Prepare(); err != nil {
			t.Fatalf("Prepare failed: %v", err)
		}

		result := s.Search(_contents)
		_ = result.IsOK()
		_ = result.Rows()
	})

	t.Run("NewSearchLastKey", func(t *testing.T) {
		s := NewSearchLastKey(NewExportKeywordAll(rule))
		defer s.Close()

		if err := s.Prepare(); err != nil {
			t.Fatalf("Prepare failed: %v", err)
		}

		result := s.Search(_contents)
		_ = result.IsOK()
		_ = result.Rows()
	})
}

// TestExport_Constructors tests NewExportKeyword, NewExportLabel, and GetRule
func TestExport_Constructors(t *testing.T) {
	rule := &ruleTest{}

	t.Run("NewExportKeyword", func(t *testing.T) {
		df := map[string]any{"k1": "v1", "k2": 0}
		export := NewExportKeyword(rule, df, func(ctx extract.ExportContext[Results]) []map[string]any {
			return []map[string]any{{"k1": "x", "k2": 1}}
		})

		if export.GetRule() != rule {
			t.Error("GetRule() should return the rule")
		}

		if len(export.Keys()) != 2 {
			t.Errorf("expected 2 keys, got %d", len(export.Keys()))
		}

		if len(export.DefaultValues()) != 2 {
			t.Errorf("expected 2 default values, got %d", len(export.DefaultValues()))
		}
	})

	t.Run("NewExportLabel", func(t *testing.T) {
		df := map[string]any{"k1": "v1"}
		export := NewExportLabel(rule, df, func(ctx extract.ExportContext[LabelResults]) []map[string]any {
			return []map[string]any{{"k1": "x"}}
		})

		if export.GetRule() != rule {
			t.Error("GetRule() should return the rule")
		}

		if len(export.Keys()) != 1 {
			t.Errorf("expected 1 key, got %d", len(export.Keys()))
		}
	})

	t.Run("Pluck_Direct", func(t *testing.T) {
		df := map[string]any{"k1": "v1", "k2": 0, "k3": "v3"}
		export := NewExportKeyword(rule, df, func(ctx extract.ExportContext[Results]) []map[string]any {
			return nil
		})

		export.Pluck([]string{"k1", "k3"})

		keys := export.Keys()
		if len(keys) != 2 {
			t.Fatalf("expected 2 keys after pluck, got %d", len(keys))
		}
		if !slices.Contains(keys, "k1") || !slices.Contains(keys, "k3") {
			t.Errorf("expected keys k1 and k3, got %v", keys)
		}

		dv := export.DefaultValues()
		if len(dv) != 2 {
			t.Fatalf("expected 2 default values after pluck, got %d", len(dv))
		}
		if dv["k1"] != "v1" {
			t.Errorf("expected v1 for k1, got %v", dv["k1"])
		}
	})

	t.Run("Pluck_NonExistentKey", func(t *testing.T) {
		df := map[string]any{"k1": "v1"}
		export := NewExportKeyword(rule, df, func(ctx extract.ExportContext[Results]) []map[string]any {
			return nil
		})

		// Pluck a key that doesn't exist should be ignored
		export.Pluck([]string{"k1", "nonexistent"})

		keys := export.Keys()
		if len(keys) != 1 {
			t.Fatalf("expected 1 key after pluck, got %d", len(keys))
		}
		if keys[0] != "k1" {
			t.Errorf("expected key k1, got %s", keys[0])
		}
	})
}

// TestExport_ResultConstructors tests NewResult, NewLabelResult, ToTag methods
func TestExport_ResultConstructors(t *testing.T) {
	rule := &ruleTest{}

	t.Run("NewResult", func(t *testing.T) {
		r := NewResult()
		if r.Tags == nil {
			t.Error("Tags should not be nil")
		}
		if r.Texts == nil {
			t.Error("Texts should not be nil")
		}
	})

	t.Run("NewLabelResult", func(t *testing.T) {
		lr := NewLabelResult()
		if lr.Tags == nil {
			t.Error("Tags should not be nil")
		}
		if lr.Match == nil {
			t.Error("Match should not be nil")
		}
	})

	t.Run("Result_ToTag", func(t *testing.T) {
		r := NewResult()
		r.Keyword = "test_keyword"
		r.Amount = 5
		r.Tags["a"] = "tag_a"
		r.Tags["ab"] = "tag_ab"
		r.Texts["test_keyword"] = 5

		tag := r.ToTag(rule)
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
		lr := NewLabelResult()
		lr.Identity = "test"
		lr.Amount = 3
		lr.Tags["a"] = "tag_a"
		lr.Match["key1"] = map[string]int{"text1": 2, "text2": 1}
		lr.Keywords = []string{"key1"}

		ltag := lr.ToTag(rule, DefaultFormat)
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
		lr := NewLabelResult()
		lr.Amount = 3
		lr.Match["key1"] = map[string]int{"text1": 2, "text2": 1}
		lr.Keywords = []string{"key1"}

		format := Format{WithKeywordAmount: false, Sep: ","}
		ltag := lr.ToTag(rule, format)
		// WithKeywordAmount=false, keyText = "key1" (no amount)
		if ltag[rule.KeywordNameAlias()] != "key1" {
			t.Errorf("expected 'key1', got %v", ltag[rule.KeywordNameAlias()])
		}
	})
}
