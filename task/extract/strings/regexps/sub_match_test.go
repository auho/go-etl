package regexps

import (
	"testing"

	"github.com/auho/go-etl/v3/internal/testutil"
)

var _rule = &ruleTest{}
var _content = []string{
	"1-2-3-12-123-1-2-3-12-123-1-2-3-12-123-1-2-3-12-123",
	"a-b-c-ab-abc-a-b-c-ab-aac-acc-a-b-c-ab-abc-a-b-c-ab-aac-acc",
}

var _expressions = []string{
	`a.{1,2}c`,
	`\b(1)\b`,
	`\b(a)\b`,
	`.*(ab).*`,
}

func TestAllSubMatch(t *testing.T) {
	amm := NewAllSubMatch(_expressions, _rule)
	err := amm.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	token := amm.Extract(_content)
	_, rets := token.Get()
	if len(rets) != 6 {
		t.Fatal()
	}

	_assertFunc := func(index int, sub string, num int) {
		if rets[index][_rule.NameAlias()] != sub || rets[index][_rule.KeywordAmountNameAlias()] != num {
			t.Fatal(index)
		}
	}

	_assertFunc(0, "1", 4)
	_assertFunc(1, "abc", 2)
	_assertFunc(2, "aac", 2)
	_assertFunc(3, "acc", 2)
	_assertFunc(4, "a", 4)
	_assertFunc(5, "ab", 1)
}

func TestSubMatchAll(t *testing.T) {
	t.Run("all", func(t *testing.T) {
		sma := NewSubMatchAll(_expressions, _rule)
		err := sma.Prepare()
		if err != nil {
			t.Fatal()
		}

		token := sma.Extract(_content)
		_, rets := token.Get()
		if len(rets) != 4 {
			t.Fatal()
		}

		_assertFunc := func(index int, sub string, num int) {
			if rets[index][_rule.NameAlias()] != sub || rets[index][_rule.KeywordAmountNameAlias()] != num {
				t.Fatal(index)
			}
		}

		_assertFunc(0, "1", 1)
		_assertFunc(1, "abc", 1)
		_assertFunc(2, "a", 1)
		_assertFunc(3, "ab", 1)
	})

	t.Run("line", func(t *testing.T) {
		sma := NewSubMatchAllLine(_expressions, _rule)
		err := sma.Prepare()
		if err != nil {
			t.Fatal()
		}

		token := sma.Extract(_content)
		_, rets := token.Get()
		if len(rets) != 1 {
			t.Fatal()
		}

		if rets[0][_rule.NameAlias()] != "1|abc|a|ab" || rets[0][_rule.KeywordAmountNameAlias()] != 4 {
			t.Fatal()
		}
	})

	t.Run("flag", func(t *testing.T) {
		sma := NewSubMatchAllFlag(_expressions, _rule)
		err := sma.Prepare()
		if err != nil {
			t.Fatal()
		}

		token := sma.Extract(_content)
		_, rets := token.Get()
		if len(rets) != 1 {
			t.Fatal()
		}

		if rets[0][_rule.NameAlias()] != 1 || rets[0][_rule.KeywordNameAlias()] != "1|abc|a|ab" {
			t.Fatal()
		}
	})
}

func TestSubMatchFirst(t *testing.T) {
	smf := NewSubMatchFirst(_expressions, _rule)
	err := smf.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	token := smf.Extract(_content)
	_, rets := token.Get()
	if len(rets) != 1 {
		t.Fatal()
	}

	if rets[0][_rule.NameAlias()] != "1" || rets[0][_rule.KeywordAmountNameAlias()] != 1 {
		t.Fatal()
	}
}

func TestSubMatch_Interface(t *testing.T) {
	t.Run("Title", func(t *testing.T) {
		sm := NewAllSubMatch(_expressions, _rule)
		expected := "SubMatch[" + _rule.Name() + "]"
		if sm.Title() != expected {
			t.Fatalf("expected Title() to be %q, got %q", expected, sm.Title())
		}
	})

	t.Run("Keys", func(t *testing.T) {
		sm := NewAllSubMatch(_expressions, _rule)
		keys := sm.Keys()
		if len(keys) != 2 {
			t.Fatalf("expected 2 keys, got %d", len(keys))
		}
		keySet := make(map[string]bool)
		for _, k := range keys {
			keySet[k] = true
		}
		if !keySet[_rule.NameAlias()] {
			t.Errorf("expected key %q", _rule.NameAlias())
		}
		if !keySet[_rule.KeywordAmountNameAlias()] {
			t.Errorf("expected key %q", _rule.KeywordAmountNameAlias())
		}
	})

	t.Run("DefaultValues", func(t *testing.T) {
		sm := NewAllSubMatch(_expressions, _rule)
		dv := sm.DefaultValues()
		if len(dv) != 2 {
			t.Fatalf("expected 2 default values, got %d", len(dv))
		}
		if _, ok := dv[_rule.NameAlias()]; !ok {
			t.Errorf("expected default value for %q", _rule.NameAlias())
		}
		if _, ok := dv[_rule.KeywordAmountNameAlias()]; !ok {
			t.Errorf("expected default value for %q", _rule.KeywordAmountNameAlias())
		}
		testutil.AssertMapCloned(t, "SubMatch", sm.DefaultValues)
	})
}

func TestSubMatch_Close(t *testing.T) {
	sm := NewAllSubMatch(_expressions, _rule)
	if err := sm.Prepare(); err != nil {
		t.Fatal(err)
	}
	if err := sm.Close(); err != nil {
		t.Fatal(err)
	}
}
