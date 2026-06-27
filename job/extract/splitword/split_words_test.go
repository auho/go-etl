package splitword

import (
	"testing"

	"github.com/auho/go-etl/v3/job/extract"
)

func TestSplitWords(t *testing.T) {
	content := "1,2,3,4,5,6,7,8,9"

	t.Run("all", func(t *testing.T) {
		s := NewSplitWords(",", NewExportAll())
		err := s.Prepare()
		if err != nil {
			t.Fatal(err)
		}

		token := s.Search([]string{content, content})
		rets := token.Rows()
		if !token.IsOK() {
			t.Fatal()
		}

		if len(rets) != 18 {
			t.Fatal()
		}
	})

	t.Run("line", func(t *testing.T) {
		format := Format{WordName: NameWord, Sep: "-"}
		e := extract.NewExporter(
			map[string]any{format.WordName: ""},
			func(ctx extract.ExportContext[Results]) []map[string]any {
				return ctx.Results.ToLine(ctx.Format.(Format))
			},
			extract.WithFormat[Results](format),
		)
		s := NewSplitWords(",", e)
		err := s.Prepare()
		if err != nil {
			t.Fatal(err)
		}

		token := s.Search([]string{content, content})
		rets := token.Rows()
		if !token.IsOK() {
			t.Fatal()
		}

		if len(rets) != 1 {
			t.Fatal()
		}

		if rets[0][NameWord] != "1-2-3-4-5-6-7-8-9-1-2-3-4-5-6-7-8-9" {
			t.Fatal()
		}
	})

	t.Run("line_default", func(t *testing.T) {
		s := NewSplitWords(",", NewExportLine())
		err := s.Prepare()
		if err != nil {
			t.Fatal(err)
		}

		token := s.Search([]string{content, content})
		rets := token.Rows()
		if !token.IsOK() {
			t.Fatal()
		}

		if len(rets) != 1 {
			t.Fatal()
		}

		if rets[0][NameWord] != "1 2 3 4 5 6 7 8 9 1 2 3 4 5 6 7 8 9" {
			t.Fatal()
		}
	})
}

func TestSplitWords_Interface(t *testing.T) {
	t.Run("Title", func(t *testing.T) {
		s := NewDefault(",")
		expected := "SplitWords[,]"
		if s.Title() != expected {
			t.Fatalf("expected Title() to be %q, got %q", expected, s.Title())
		}
	})

	t.Run("NewExport", func(t *testing.T) {
		s := NewDefault(",")
		if s.NewExport() == nil {
			t.Fatal("NewExport() returned nil")
		}
	})

	t.Run("Keys", func(t *testing.T) {
		e := NewExportAll()
		keys := e.Keys()
		if len(keys) != 1 {
			t.Fatalf("expected 1 key, got %d", len(keys))
		}
		if keys[0] != NameWord {
			t.Errorf("expected key %q, got %q", NameWord, keys[0])
		}
	})

	t.Run("DefaultValues", func(t *testing.T) {
		e := NewExportAll()
		dv := e.DefaultValues()
		if len(dv) != 1 {
			t.Fatalf("expected 1 default value, got %d", len(dv))
		}
		if _, ok := dv[NameWord]; !ok {
			t.Errorf("expected default value for %q", NameWord)
		}
	})

	t.Run("Prepare and Close", func(t *testing.T) {
		s := NewDefault(",")
		if err := s.Prepare(); err != nil {
			t.Fatal(err)
		}
		if err := s.Close(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("WithFormat", func(t *testing.T) {
		customFormat := Format{
			WordName: "custom_word",
			Sep:      "-",
		}
		e := extract.NewExporter(
			map[string]any{customFormat.WordName: ""},
			func(ctx extract.ExportContext[Results]) []map[string]any {
				return ctx.Results.ToLine(ctx.Format.(Format))
			},
			extract.WithFormat[Results](customFormat),
		)
		s := NewSplitWords(",", e)
		if err := s.Prepare(); err != nil {
			t.Fatal(err)
		}
		defer s.Close()

		token := s.Search([]string{"a,b,c"})
		rets := token.Rows()
		if len(rets) != 1 {
			t.Fatalf("expected 1 row, got %d", len(rets))
		}
		if rets[0]["custom_word"] != "a-b-c" {
			t.Errorf("expected %q, got %v", "a-b-c", rets[0]["custom_word"])
		}
	})

	t.Run("FormatCheck", func(t *testing.T) {
		f := Format{}
		f.check()
		if f.WordName != NameWord {
			t.Errorf("expected WordName %q, got %q", NameWord, f.WordName)
		}
		if f.Sep != " " {
			t.Errorf("expected Sep %q, got %q", " ", f.Sep)
		}
	})
}
