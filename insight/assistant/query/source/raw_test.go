package source

import (
	"testing"

	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
	"github.com/auho/go-etl/v3/insight/assistant/sqlbuilder/dml"
	simpledb "github.com/auho/go-simple-db/v3"
)

// mockRaw implements assistant.Raw for testing.
var _ assistant.Raw = (*mockRaw)(nil)

type mockRaw struct {
	db        *simpledb.SimpleDB
	name      string
	tableName string
	dmlTable  *dml.Table
}

func (m *mockRaw) DB() *simpledb.SimpleDB     { return m.db }
func (m *mockRaw) Name() string                { return m.name }
func (m *mockRaw) TableName() string           { return m.tableName }
func (m *mockRaw) ExecCommand(*schema.Command) {}
func (m *mockRaw) DMLTable() *dml.Table        { return m.dmlTable }

// --- unit tests ---

func TestNewRaw(t *testing.T) {
	mock := &mockRaw{name: "test", tableName: "test_table"}
	r := NewRaw("test", mock)

	if r.Name != "test" {
		t.Fatalf("expect[test] != actual[%s]", r.Name)
	}
	if r.Raw != mock {
		t.Fatal("Raw field not set correctly")
	}
	if r.Base.Name != "test" {
		t.Fatalf("expect Base.Name[test] != actual[%s]", r.Base.Name)
	}
}

// --- error path tests ---

func TestRawDatasetTableNotFound(t *testing.T) {
	skipIfNoDB(t)

	mock := &mockRaw{
		db:        _simpleDB,
		name:      "raw_err",
		tableName: "nonexistent_table",
		dmlTable:  dml.NewTable("nonexistent_table"),
	}

	r := NewRaw("raw_err", mock)
	_, err := r.Dataset()
	if err == nil {
		t.Fatal("expect error for non-existent table")
	}
}

// --- integration tests ---

func TestRawDataset(t *testing.T) {
	skipIfNoDB(t)

	mock := &mockRaw{
		db:        _simpleDB,
		name:      "raw",
		tableName: _testTable,
		dmlTable:  dml.NewTable(_testTable),
	}

	r := NewRaw("raw", mock)
	ds, err := r.Dataset()
	if err != nil {
		t.Fatal(err)
	}

	if ds.Name != "raw" {
		t.Fatalf("expect[raw] != actual[%s]", ds.Name)
	}
	// test_source table has 4 columns: id, name, category, value
	if len(ds.Titles) != 4 {
		t.Fatalf("expect[4] != actual[%d]", len(ds.Titles))
	}
	if len(ds.Sets) != 1 {
		t.Fatalf("expect[1] != actual[%d]", len(ds.Sets))
	}
	if ds.Sets[0].Amount != 4 {
		t.Fatalf("expect[4] != actual[%d]", ds.Sets[0].Amount)
	}
}

func TestRawDatasetWithWhere(t *testing.T) {
	skipIfNoDB(t)

	mock := &mockRaw{
		db:        _simpleDB,
		name:      "raw_cat1",
		tableName: _testTable,
		dmlTable:  dml.NewTable(_testTable).Where("category = 'cat1'"),
	}

	r := NewRaw("raw_cat1", mock)
	ds, err := r.Dataset()
	if err != nil {
		t.Fatal(err)
	}

	if len(ds.Sets) != 1 {
		t.Fatalf("expect[1] != actual[%d]", len(ds.Sets))
	}
	// cat1 has 2 rows
	if ds.Sets[0].Amount != 2 {
		t.Fatalf("expect[2] != actual[%d]", ds.Sets[0].Amount)
	}
}

func TestRawDatasetEmptyResult(t *testing.T) {
	skipIfNoDB(t)

	mock := &mockRaw{
		db:        _simpleDB,
		name:      "raw_empty",
		tableName: _testTable,
		dmlTable:  dml.NewTable(_testTable).Where("category = 'notexist'"),
	}

	r := NewRaw("raw_empty", mock)
	ds, err := r.Dataset()
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

func TestRawDatasetQueryError(t *testing.T) {
	skipIfNoDB(t)

	// tableName is valid so GetTableColumns succeeds,
	// but dmlTable points to a non-existent table so the SQL query fails
	mock := &mockRaw{
		db:        _simpleDB,
		name:      "raw_query_err",
		tableName: _testTable,
		dmlTable:  dml.NewTable("nonexistent_table"),
	}

	r := NewRaw("raw_query_err", mock)
	_, err := r.Dataset()
	if err == nil {
		t.Fatal("expect error for query on non-existent table")
	}
}
