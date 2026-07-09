package source

import (
	"testing"

	"github.com/auho/go-etl/v3/insight/assistant/sqlbuilder/dml"
)

func TestRowsStackSourceDataset(t *testing.T) {
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
