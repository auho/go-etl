package containment

import (
	"math/rand"
	"testing"

	"github.com/auho/go-etl/v3/internal/testutil"
)

var _rule = &ruleTest{}
var _content = []string{
	"1-2-3-12-123",
	"a-b-c-ab-abc",
}

func TestContainsAll(t *testing.T) {
	subs := []string{"1", "2", "12", "ab"}

	t.Run("all", func(t *testing.T) {
		c := NewContainsAll(subs, _rule)
		err := c.Prepare()
		if err != nil {
			t.Fatal(err)
		}

		token := c.Extract(_content)
		if !token.IsOK() {
			t.Fatal()
		}

		_, rets := token.Get()
		if len(rets) != 4 {
			t.Fatal()
		}

		_assertFunc := func(index int, sub string, num int) {
			if rets[index][_rule.NameAlias()] != sub || rets[index][_rule.KeywordAmountNameAlias()] != num {
				t.Fatal(index)
			}
		}
		_assertFunc(0, "1", 3)
		_assertFunc(1, "2", 3)
		_assertFunc(2, "12", 2)
		_assertFunc(3, "ab", 2)
	})

	t.Run("line", func(t *testing.T) {
		c := NewContainsAllLine(subs, _rule)
		err := c.Prepare()
		if err != nil {
			t.Fatal(err)
		}

		token := c.Extract(_content)
		if !token.IsOK() {
			t.Fatal()
		}

		_, rets := token.Get()
		if len(rets) != 1 {
			t.Fatal()
		}

		if rets[0][_rule.NameAlias()] != "1|2|12|ab" || rets[0][_rule.KeywordAmountNameAlias()] != 10 {
			t.Fatal()
		}
	})

	t.Run("flag", func(t *testing.T) {
		c := NewContainsAllFlag(subs, _rule)
		err := c.Prepare()
		if err != nil {
			t.Fatal(err)
		}

		token := c.Extract(_content)
		if !token.IsOK() {
			t.Fatal()
		}

		_, rets := token.Get()
		if len(rets) != 1 {
			t.Fatal()
		}

		if rets[0][_rule.NameAlias()] != 1 || rets[0][_rule.KeywordNameAlias()] != "1|2|12|ab" {
			t.Fatal()
		}
	})
}

func TestContainsAny(t *testing.T) {
	ss := []string{"1", "a"}

	for i := 0; i < 5; i++ {
		t.Run("any", func(t *testing.T) {
			rand.Shuffle(len(ss), func(i, j int) {
				ss[i], ss[j] = ss[j], ss[i]
			})

			c := NewContainsFirst(ss, _rule)
			err := c.Prepare()
			if err != nil {
				t.Fatal(err)
			}

			token := c.Extract(_content)
			if !token.IsOK() {
				t.Fatal()
			}

			_, rets := token.Get()
			if len(rets) != 1 {
				t.Fatal()
			}

			if rets[0][_rule.NameAlias()] != "1" || rets[0][_rule.KeywordAmountNameAlias()] != 3 {
				t.Fatal()
			}
		})

	}
}

func TestContains_Interface(t *testing.T) {
	t.Run("Title", func(t *testing.T) {
		c := NewContainsAll([]string{"1"}, _rule)
		expected := "Contains[" + _rule.Name() + "]"
		if c.Title() != expected {
			t.Fatalf("expected Title() to be %q, got %q", expected, c.Title())
		}
	})

	t.Run("Keys", func(t *testing.T) {
		c := NewContainsAll([]string{"1"}, _rule)
		keys := c.Keys()
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
		c := NewContainsAll([]string{"1"}, _rule)
		dv := c.DefaultValues()
		if len(dv) != 2 {
			t.Fatalf("expected 2 default values, got %d", len(dv))
		}
		if _, ok := dv[_rule.NameAlias()]; !ok {
			t.Errorf("expected default value for %q", _rule.NameAlias())
		}
		if _, ok := dv[_rule.KeywordAmountNameAlias()]; !ok {
			t.Errorf("expected default value for %q", _rule.KeywordAmountNameAlias())
		}
		testutil.AssertMapCloned(t, "Contains", c.DefaultValues)
	})
}

func TestContains_Close(t *testing.T) {
	c := NewContainsAll([]string{"1"}, _rule)
	if err := c.Prepare(); err != nil {
		t.Fatal(err)
	}
	if err := c.Close(); err != nil {
		t.Fatal(err)
	}
}
