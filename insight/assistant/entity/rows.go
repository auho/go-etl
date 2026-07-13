package entity

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ assistant.Entity = (*Rows)(nil)

// Rows represents an intermediate row-level data table in the ETL pipeline.
// It is a named table with an id column, used for storing and processing row data
// between the raw source stage and the final data stage.
type Rows struct {
	base             // embedded base for command hook and DB connection
	extra            // embedded extra for DDL/DML operations
	name      string // entity name
	idName    string // name of the primary key column
	tableName string // actual database table name
}

// NewRowsCustomTable creates a new Rows entity with a custom table name.
func NewRowsCustomTable(name, tableName, idName string, db *simpledb.SimpleDB) *Rows {
	r := &Rows{
		base:      base{db: db},
		name:      name,
		idName:    idName,
		tableName: tableName,
	}
	r.extra = extra{
		raw: r,
	}

	return r
}

// NewRows creates a new Rows entity where the table name equals the entity name.
func NewRows(name, idName string, db *simpledb.SimpleDB) *Rows {
	return NewRowsCustomTable(name, name, idName, db)
}

// Name returns the entity name.
func (r *Rows) Name() string {
	return r.name
}

// IDName returns the name of the primary key column.
func (r *Rows) IDName() string {
	return r.idName
}

// TableName returns the actual database table name.
func (r *Rows) TableName() string {
	return r.tableName
}

// WithCommand sets the command hook and returns the Rows instance for chaining.
func (r *Rows) WithCommand(fn func(command *schema.Command)) *Rows {
	r.withCommand(fn)

	return r
}

// ToData converts this Rows into a Data entity with the same name and id column.
func (r *Rows) ToData() *Data {
	return NewData(r.name, r.idName, r.db)
}

// ToRaw converts this Rows into a Raw entity with the same name.
func (r *Rows) ToRaw() *Raw {
	return NewRaw(r.name, r.db)
}

// Clone creates a new Rows entity with the given name, preserving the id column and command hook.
func (r *Rows) Clone(name string) *Rows {
	return NewRows(name, r.idName, r.db).WithCommand(r.commandFunc)
}

// CloneSuffix creates a new Rows entity by appending the given suffixes to the current table name.
// Hyphens in suffixes are replaced with underscores.
func (r *Rows) CloneSuffix(suffix ...string) *Rows {
	var ns []string
	for _, _s := range suffix {
		ns = append(ns, strings.ReplaceAll(_s, "-", "_"))
	}

	return r.Clone(strings.Join(append([]string{r.TableName()}, ns...), "_"))
}

// ToDeletedRows creates a Rows entity for tracking deleted rows,
// with the table name prefixed by "deleted_".
func (r *Rows) ToDeletedRows() *Rows {
	return NewRows(fmt.Sprintf("%s_%s", NameDeleted, r.name), r.idName, r.db)
}
