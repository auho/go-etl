package dml

import (
	simpledb "github.com/auho/go-simple-db/v2"
)

type Tabler interface {
	Manipulator
	GetSelectFields() []string
}

type statement interface {
	Query() string
	InsertQuery(name string) string
	InsertWithFieldsQuery(name string, fields []string) string
	UpdateQuery() string
	DeleteQuery() string
}

type Manipulator interface {
	SQL() string
	InsertSQL(name string) string
	Insert(name string, db *simpledb.SimpleDB) (string, error)
	InsertWithFieldsSQL(name string, fields []string) string
	InsertWithField(name string, fields []string, db *simpledb.SimpleDB) (string, error)
	UpdateSQL() string
	Update(db *simpledb.SimpleDB) (string, error)
	DeleteSQL() string
	Delete(db *simpledb.SimpleDB) (string, error)
}
