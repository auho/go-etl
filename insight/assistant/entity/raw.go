package entity

import (
	"github.com/auho/go-etl/v3/insight/assistant"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ assistant.Raw = (*Raw)(nil)

// Raw represents a source data table entity.
// It maps directly to an existing database table without any naming transformation.
// This is the most basic entity type, used as the entry point for the data pipeline.
type Raw struct {
	base  // embedded base for command hook and DB connection
	extra // embedded extra for DDL/DML operations
	name  string // the actual database table name
}

// NewRaw creates a new Raw entity with the given table name and database connection.
func NewRaw(name string, db *simpledb.SimpleDB) *Raw {
	r := &Raw{}
	r.name = name
	r.db = db
	r.extra = extra{
		base: r,
	}

	return r
}

// DB returns the database connection.
func (r *Raw) DB() *simpledb.SimpleDB {
	return r.db
}

// Name returns the entity's name, which is the raw table name.
func (r *Raw) Name() string {
	return r.name
}

// TableName returns the actual database table name, which is the same as Name for Raw.
func (r *Raw) TableName() string {
	return r.name
}
