package splitword

import (
	"testing"

	"github.com/auho/go-etl/v3/internal/testutil"
)

func TestSplitWords(t *testing.T) {
	content := "1,2,3,4,5,6,7,8,9"

	t.Run("all", func(t *testing.T) {
		s := NewSplitWordsAll(",")
		err := s.Prepare()
		if err != nil {
			t.Fatal(err)
		}

		token := s.Extract([]string{content, content})
		if !token.IsOK() {
			t.Fatal()
		}
		_, rets := token.Get()

		if len(rets) != 18 {
			t.Fatalf("expected 18 rows, got %d", len(rets))
		}
	})

	t.Run("line", func(t *testing.T) {
		s := NewSplitWordsLine(",")
		s.format = format{wordName: nameWord, sep: "-"}
		err := s.Prepare()
		if err != nil {
			t.Fatal(err)
		}

		token := s.Extract([]string{content, content})
		if !token.IsOK() {
			t.Fatal()
		}
		_, rets := token.Get()

		if len(rets) != 1 {
			t.Fatalf("expected 1 row, got %d", len(rets))
		}
		if rets[0][nameWord] != "1-2-3-4-5-6-7-8-9-1-2-3-4-5-6-7-8-9" {
			t.Fatalf("got %v", rets[0][nameWord])
		}
	})

	t.Run("line_default", func(t *testing.T) {
		s := NewSplitWordsLine(",")
		err := s.Prepare()
		if err != nil {
			t.Fatal(err)
		}

		token := s.Extract([]string{content, content})
		if !token.IsOK() {
			t.Fatal()
		}
		_, rets := token.Get()

		if len(rets) != 1 {
			t.Fatalf("expected 1 row, got %d", len(rets))
		}
		if rets[0][nameWord] != "1 2 3 4 5 6 7 8 9 1 2 3 4 5 6 7 8 9" {
			t.Fatalf("got %v", rets[0][nameWord])
		}
	})
}

func TestSplitWords_Interface(t *testing.T) {
	t.Run("Title", func(t *testing.T) {
		s := NewSplitWordsAll(",")
		expected := "SplitWords[,]"
		if s.Title() != expected {
			t.Fatalf("expected Title() to be %q, got %q", expected, s.Title())
		}
	})

	t.Run("Keys", func(t *testing.T) {
		s := NewSplitWordsAll(",")
		if err := s.Prepare(); err != nil {
			t.Fatal(err)
		}
		keys := s.Keys()
		if len(keys) != 1 {
			t.Fatalf("expected 1 key, got %d", len(keys))
		}
		if keys[0] != nameWord {
			t.Errorf("expected key %q, got %q", nameWord, keys[0])
		}
	})

	t.Run("DefaultValues", func(t *testing.T) {
		s := NewSplitWordsAll(",")
		if err := s.Prepare(); err != nil {
			t.Fatal(err)
		}
		dv := s.DefaultValues()
		if len(dv) != 1 {
			t.Fatalf("expected 1 default value, got %d", len(dv))
		}
		if _, ok := dv[nameWord]; !ok {
			t.Errorf("expected default value for %q", nameWord)
		}
		testutil.AssertMapCloned(t, "SplitWords", s.DefaultValues)
	})

	t.Run("Prepare and Close", func(t *testing.T) {
		s := NewSplitWordsAll(",")
		if err := s.Prepare(); err != nil {
			t.Fatal(err)
		}
		if err := s.Close(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("FormatCheck", func(t *testing.T) {
		f := format{}
		f.check()
		if f.wordName != nameWord {
			t.Errorf("expected wordName %q, got %q", nameWord, f.wordName)
		}
		if f.sep != " " {
			t.Errorf("expected sep %q, got %q", " ", f.sep)
		}
	})
}
