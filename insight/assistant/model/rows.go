package model

import (
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ Rowsor = (*Rows)(nil)

type Rows struct {
	tableName string
	db        *simpleDb.SimpleDB
}

func (r *Rows) GetDB() *simpleDb.SimpleDB {
	return r.db
}

func (r *Rows) TableName() string {
	return r.tableName
}
