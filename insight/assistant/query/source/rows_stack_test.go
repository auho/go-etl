package source

import (
	"testing"

	"github.com/auho/go-etl/v3/insight/assistant/sqlbuilder/dml"
)

func TestRowsStackDatasetEmpty(t *testing.T) {
	s := NewRowsStack("empty")
	_, err := s.Dataset()
	if err == nil {
		t.Fatal("expect error for empty rowsSources")
	}
}

func TestRowsStackDataset(t *testing.T) {
	skipIfNoDB(t)

	s := NewRowsStack("stack",
		Base{
			Name:  "cat1",
			Table: dml.NewTable(_testTable).Select([]string{"name", "value"}).Where("category = 'cat1'"),
			DB:    _simpleDB,
		},
		Base{
			Name:  "cat2",
			Table: dml.NewTable(_testTable).Select([]string{"name", "value"}).Where("category = 'cat2'"),
			DB:    _simpleDB,
		},
	)

	ds, err := s.Dataset()
	if err != nil {
		t.Fatal(err)
	}

	if len(ds.Sets) != 2 {
		t.Fatalf("expect[2] != actual[%d]", len(ds.Sets))
	}
}

func TestRowsStackDatasetSubSourceError(t *testing.T) {
	skipIfNoDB(t)

	s := NewRowsStack("stack_err",
		Base{
			Name:  "valid",
			Table: dml.NewTable(_testTable).Select([]string{"name"}),
			DB:    _simpleDB,
		},
		Base{
			Name:  "invalid",
			Table: dml.NewTable("nonexistent_table").Select([]string{"name"}),
			DB:    _simpleDB,
		},
	)

	_, err := s.Dataset()
	if err == nil {
		t.Fatal("expect error for sub-source with non-existent table")
	}
}
