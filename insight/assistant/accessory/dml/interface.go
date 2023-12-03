package dml

import (
	"github.com/auho/go-simple-db/v2"
)

type Tabler interface {
	Manipulationor
	GetSelectFields() []string
}

type statement interface {
	Query() string
	InsertQuery(name string) string
	InsertWithFieldsQuery(name string, fields []string) string
	UpdateQuery() string
	DeleteQuery() string
}

type Manipulationor interface {
	Sql() string
	InsertSql(name string) string
	Insert(name string, db *go_simple_db.SimpleDB) (string, error)
	InsertWithFieldsSql(name string, fields []string) string
	InsertWithField(name string, fields []string, db *go_simple_db.SimpleDB) (string, error)
	UpdateSql() string
	Update(db *go_simple_db.SimpleDB) (string, error)
	DeleteSql() string
	Delete(db *go_simple_db.SimpleDB) (string, error)
}
