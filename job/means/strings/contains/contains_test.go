package contains

import (
	"math/rand"
	"testing"
)

var _rule = &ruleTest{}
var _content = []string{
	"1-2-3-12-123",
	"a-b-c-ab-abc",
}

func TestContainsAll(t *testing.T) {
	subs := []string{"1", "2", "12", "ab"}

	t.Run("all", func(t *testing.T) {
		c := NewContainsAll(_rule, subs, NewExportAll(_rule))
		err := c.Prepare()
		if err != nil {
			t.Fatal(err)
		}

		token := c.Do(_content)
		if !token.IsOk() {
			t.Fatal()
		}

		rets := token.ToToken()
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
		c := NewContainsAll(_rule, subs, NewExportLine(_rule))
		err := c.Prepare()
		if err != nil {
			t.Fatal(err)
		}

		token := c.Do(_content)
		if !token.IsOk() {
			t.Fatal()
		}

		rets := token.ToToken()
		if len(rets) != 1 {
			t.Fatal()
		}

		if rets[0][_rule.NameAlias()] != "1|2|12|ab" || rets[0][_rule.KeywordAmountNameAlias()] != 10 {
			t.Fatal()
		}
	})

	t.Run("flag", func(t *testing.T) {
		c := NewContainsAll(_rule, subs, NewExportFlag(_rule))
		err := c.Prepare()
		if err != nil {
			t.Fatal(err)
		}

		token := c.Do(_content)
		if !token.IsOk() {
			t.Fatal()
		}

		rets := token.ToToken()
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

			c := NewContainsAny(_rule, ss, NewExportAll(_rule))
			err := c.Prepare()
			if err != nil {
				t.Fatal(err)
			}

			token := c.Do(_content)
			if !token.IsOk() {
				t.Fatal()
			}

			rets := token.ToToken()
			if len(rets) != 1 {
				t.Fatal()
			}

			if rets[0][_rule.NameAlias()] != "1" || rets[0][_rule.KeywordAmountNameAlias()] != 3 {
				t.Fatal()
			}
		})

	}
}
