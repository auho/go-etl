package entity

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ assistant.Entity = (*Rows)(nil)

type Rows struct {
	model
	extra
	name      string
	idName    string
	tableName string
}

func NewRowsCustomTable(name, tableName, idName string, db *simpledb.SimpleDB) *Rows {
	r := &Rows{}
	r.name = name
	r.idName = idName
	r.tableName = tableName
	r.db = db
	r.extra = extra{
		model: r,
	}

	return r
}

func NewRows(name, idName string, db *simpledb.SimpleDB) *Rows {
	return NewRowsCustomTable(name, name, idName, db)
}

func (r *Rows) DB() *simpledb.SimpleDB {
	return r.db
}

func (r *Rows) Name() string {
	return r.name
}

func (r *Rows) IDName() string {
	return r.idName
}

func (r *Rows) TableName() string {
	return r.tableName
}

func (r *Rows) WithCommand(fn func(command *schema.Command)) *Rows {
	r.withCommand(fn)

	return r
}

func (r *Rows) ToData() *Data {
	return NewData(r.name, r.idName, r.db)
}

func (r *Rows) ToRaw() *Raw {
	return NewRaw(r.name, r.db)
}

func (r *Rows) Clone(name string) *Rows {
	return NewRows(name, r.idName, r.db).WithCommand(r.commandFunc)
}

func (r *Rows) CloneSuffix(suffix ...string) *Rows {
	var ns []string
	for _, _s := range suffix {
		ns = append(ns, strings.ReplaceAll(_s, "-", "_"))
	}

	return r.Clone(strings.Join(append([]string{r.TableName()}, ns...), "_"))
}

func (r *Rows) ToDeletedRows() *Rows {
	return NewRows(fmt.Sprintf("%s_%s", NameDeleted, r.name), r.idName, r.db)
}
