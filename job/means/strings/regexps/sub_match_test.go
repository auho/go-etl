package regexps

import (
	"testing"
)

var _rule = &ruleTest{}
var _content = []string{
	"1-2-3-12-123-1-2-3-12-123-1-2-3-12-123-1-2-3-12-123",
	"a-b-c-ab-abc-a-b-c-ab-aac-acc-a-b-c-ab-abc-a-b-c-ab-aac-acc",
}
var _expressons = []string{
	`a.{1,2}c`,
	`\b(1)\b`,
	`\b(a)\b`,
	`.*(ab).*`,
}

func TestAllSubMatch(t *testing.T) {
	amm := NewAllSubMatch(_rule, _expressons, NewExportAll(_rule))
	err := amm.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	token := amm.Do(_content)
	rets := token.ToToken()
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
		sma := NewSubMatchAll(_rule, _expressons, NewExportAll(_rule))
		err := sma.Prepare()
		if err != nil {
			t.Fatal()
		}

		token := sma.Do(_content)
		rets := token.ToToken()
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
		sma := NewSubMatchAll(_rule, _expressons, NewExportLine(_rule))
		err := sma.Prepare()
		if err != nil {
			t.Fatal()
		}

		token := sma.Do(_content)
		rets := token.ToToken()
		if len(rets) != 1 {
			t.Fatal()
		}

		if rets[0][_rule.NameAlias()] != "1|abc|a|ab" || rets[0][_rule.KeywordAmountNameAlias()] != 4 {
			t.Fatal()
		}
	})

	t.Run("flag", func(t *testing.T) {
		sma := NewSubMatchAll(_rule, _expressons, NewExportFlag(_rule))
		err := sma.Prepare()
		if err != nil {
			t.Fatal()
		}

		token := sma.Do(_content)
		rets := token.ToToken()
		if len(rets) != 1 {
			t.Fatal()
		}

		if rets[0][_rule.NameAlias()] != 1 || rets[0][_rule.KeywordNameAlias()] != "1|abc|a|ab" {
			t.Fatal()
		}
	})
}

func TestSubMatchFirst(t *testing.T) {
	smf := NewSubMatchFirst(_rule, _expressons, NewExportAll(_rule))
	err := smf.Prepare()
	if err != nil {
		t.Fatal(err)
	}

	token := smf.Do(_content)
	rets := token.ToToken()
	if len(rets) != 1 {
		t.Fatal()
	}

	if rets[0][_rule.NameAlias()] != "1" || rets[0][_rule.KeywordAmountNameAlias()] != 1 {
		t.Fatal()
	}
}
