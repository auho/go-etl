package model

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ assistant.Rowsor = (*Rows)(nil)

type Rows struct {
	model
	extra
	name      string
	idName    string
	tableName string
}

func NewRowsFake(name, tableName, idName string, db *simpleDb.SimpleDB) *Rows {
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

func NewRows(name, idName string, db *simpleDb.SimpleDB) *Rows {
	return NewRowsFake(name, name, idName, db)
}

func (r *Rows) GetDB() *simpleDb.SimpleDB {
	return r.db
}

func (r *Rows) GetName() string {
	return r.name
}

func (r *Rows) GetIdName() string {
	return r.idName
}

func (r *Rows) TableName() string {
	return r.tableName
}

func (r *Rows) WithCommand(fn func(command *tablestructure.Command)) *Rows {
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
	return NewRows(name, r.idName, r.db).WithCommand(r.commandFun)
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
