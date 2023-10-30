package table

import (
	"github.com/auho/go-etl/v2/insight/model/ddl/command/mysql"
	simpleDb "github.com/auho/go-simple-db/v2"
)

type table struct {
	table *mysql.Table
}

func (t *table) initTable(name string) {
	t.table = &mysql.Table{}
	t.table.SetName(name).SetEngineMyISAM()
}

func (t *table) GetTable() *mysql.Table {
	return t.table
}

func (t *table) AddPkBigInt(name string) {
	t.table.AddPkBigInt(name, 20)
}

func (t *table) AddPkInt(name string) {
	t.table.AddPKInt(name, 11)
}

func (t *table) AddPkString(name string, length int) {
	t.table.AddPkVarchar(name, length)
}

func (t *table) AddKeyBigInt(name string) {
	t.table.AddBigInt(name, 20, 0, false)
	t.table.AddKey(name, 0)
}

func (t *table) AddKeyInt(name string) {
	t.table.AddInt(name, 11, 0, false)
	t.table.AddKey(name, 0)
}

func (t *table) AddKeyString(name string, length int, size int) {
	t.table.AddVarchar(name, length, "")
	t.table.AddKey(name, size)
}

func (t *table) AddUniqueBigInt(name string) {
	t.table.AddBigInt(name, 20, 0, false)
	t.table.AddUniqueKey(name)
}

func (t *table) AddUniqueInt(name string) {
	t.table.AddInt(name, 11, 0, false)
	t.table.AddUniqueKey(name)
}

func (t *table) AddUniqueString(name string, length int) {
	t.table.AddVarchar(name, length, "")
	t.table.AddUniqueKey(name)
}

func (t *table) AddInt(name string) {
	t.table.AddInt(name, 11, 0, false)
}

func (t *table) AddBigInt(name string) {
	t.table.AddBigInt(name, 11, 0, false)
}

func (t *table) AddString(name string) {
	t.table.AddVarchar(name, 100, "")
}

func (t *table) AddStringWithLength(name string, length int) {
	t.table.AddVarchar(name, length, "")
}

func (t *table) Build(db *simpleDb.SimpleDB) error {
	sql := t.table.SqlForCreate()

	return db.Exec(sql).Error
}
