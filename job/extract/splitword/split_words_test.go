package splitword

import (
	"testing"
)

func TestMeans(t *testing.T) {
	content := "1,2,3,4,5,6,7,8,9"

	t.Run("all", func(t *testing.T) {
		s := NewSplitWords(",", NewExportAll())
		err := s.Prepare()
		if err != nil {
			t.Fatal(err)
		}

		token := s.Do([]string{content, content})
		rets := token.ToToken()
		if !token.IsOk() {
			t.Fatal()
		}

		if len(rets) != 18 {
			t.Fatal()
		}
	})

	t.Run("line", func(t *testing.T) {
		s := NewSplitWords(",", NewExportLine().WithFormat(Format{Sep: "-"}))
		err := s.Prepare()
		if err != nil {
			t.Fatal(err)
		}

		token := s.Do([]string{content, content})
		rets := token.ToToken()
		if !token.IsOk() {
			t.Fatal()
		}

		if len(rets) != 1 {
			t.Fatal()
		}

		if rets[0][NameWord] != "1-2-3-4-5-6-7-8-9-1-2-3-4-5-6-7-8-9" {
			t.Fatal()
		}
	})
}
