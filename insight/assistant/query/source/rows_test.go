package source

import (
	"testing"

	"github.com/auho/go-etl/v3/insight/assistant/sqlbuilder/dml"
)

// --- error path tests ---

func TestRowsDatasetTableNotFound(t *testing.T) {
	skipIfNoDB(t)

	s := NewRows(Base{
		Name:  "rows_err",
		Table: dml.NewTable("nonexistent_table").Select([]string{"name"}),
		DB:    _simpleDB,
	})

	_, err := s.Dataset()
	if err == nil {
		t.Fatal("expect error for non-existent table")
	}
}

// --- integration tests ---

func TestRowsDataset(t *testing.T) {
	skipIfNoDB(t)

	s := NewRows(Base{
		Name:  "rows",
		Table: dml.NewTable(_testTable).Select([]string{"name", "category", "value"}),
		DB:    _simpleDB,
	})

	ds, err := s.Dataset()
	if err != nil {
		t.Fatal(err)
	}

	if ds.Name != "rows" {
		t.Fatalf("expect[rows] != actual[%s]", ds.Name)
	}
	if len(ds.Titles) != 3 {
		t.Fatalf("expect[3] != actual[%d]", len(ds.Titles))
	}
	if len(ds.Sets) != 1 {
		t.Fatalf("expect[1] != actual[%d]", len(ds.Sets))
	}
	if ds.Sets[0].Amount != 4 {
		t.Fatalf("expect[4] != actual[%d]", ds.Sets[0].Amount)
	}
}

func TestRowsDatasetWithWhere(t *testing.T) {
	skipIfNoDB(t)

	s := NewRows(Base{
		Name:  "rows_cat1",
		Table: dml.NewTable(_testTable).Select([]string{"name", "value"}).Where("category = 'cat1'"),
		DB:    _simpleDB,
	})

	ds, err := s.Dataset()
	if err != nil {
		t.Fatal(err)
	}

	if len(ds.Sets) != 1 {
		t.Fatalf("expect[1] != actual[%d]", len(ds.Sets))
	}
	if ds.Sets[0].Amount != 2 {
		t.Fatalf("expect[2] != actual[%d]", ds.Sets[0].Amount)
	}
}

func TestRowsDatasetEmptyResult(t *testing.T) {
	skipIfNoDB(t)

	s := NewRows(Base{
		Name:  "rows_empty",
		Table: dml.NewTable(_testTable).Select([]string{"name"}).Where("category = 'notexist'"),
		DB:    _simpleDB,
	})

	ds, err := s.Dataset()
	if err != nil {
		t.Fatal(err)
	}

	if len(ds.Sets) != 1 {
		t.Fatalf("expect[1] != actual[%d]", len(ds.Sets))
	}
	if ds.Sets[0].Amount != 0 {
		t.Fatalf("expect[0] != actual[%d]", ds.Sets[0].Amount)
	}
}
