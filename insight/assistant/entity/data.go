package entity

import (
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ assistant.Entity = (*Data)(nil)

// Data represents a processed data table in the ETL pipeline.
// Its actual database table name is prefixed with "data_" (e.g., "data_orders").
// Data is the output of the row-level processing stage and serves as input
// for rule-based analysis and tagging.
type Data struct {
	base   // embedded base for command hook and DB connection
	extra  // embedded extra for DDL/DML operations
	name   string // entity name (used to derive the table name)
	idName string // name of the primary key column
}

// NewData creates a new Data entity with the given name, id column, and database connection.
func NewData(name string, idName string, db *simpledb.SimpleDB) *Data {
	d := &Data{}
	d.name = name
	d.idName = idName
	d.db = db
	d.extra = extra{
		base: d,
	}

	return d
}

// DB returns the database connection.
func (d *Data) DB() *simpledb.SimpleDB {
	return d.db
}

// Name returns the entity name.
func (d *Data) Name() string {
	return d.name
}

// IDName returns the name of the primary key column.
func (d *Data) IDName() string {
	return d.idName
}

// TableName returns the database table name, prefixed with "data_".
func (d *Data) TableName() string {
	return fmt.Sprintf("%s_%s", NameData, d.name)
}

// WithCommand sets the command hook and returns the Data instance for chaining.
func (d *Data) WithCommand(fn func(command *schema.Command)) *Data {
	d.withCommand(fn)

	return d
}

// ToRows converts this Data into a Rows entity with the same name and id column.
func (d *Data) ToRows() *Rows {
	return NewRows(d.name, d.idName, d.db)
}

// ToRaw converts this Data into a Raw entity with the same name.
func (d *Data) ToRaw() *Raw {
	return NewRaw(d.name, d.db)
}
