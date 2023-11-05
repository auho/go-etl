package table

import (
	"github.com/auho/go-etl/v2/insight/assistant/accessory/ddl/command/mysql"
)

type command struct {
	table *mysql.Table
}

func (c *command) TableName() string {
	return c.table.GetName()
}

func (c *command) Sql() string {
	return c.table.SqlForCreate()
}

func (c *command) AddPkBigInt(name string) {
	c.table.AddPkBigInt(name, 20)
}

func (c *command) AddPkInt(name string) {
	c.table.AddPKInt(name, 11)
}

func (c *command) AddPkString(name string, length int) {
	c.table.AddPkVarchar(name, length)
}

func (c *command) AddKeyBigInt(name string) {
	c.table.AddBigInt(name, 20, 0, false)
	c.table.AddKey(name, 0)
}

func (c *command) AddKeyInt(name string) {
	c.table.AddInt(name, 11, 0, false)
	c.table.AddKey(name, 0)
}

func (c *command) AddKeyString(name string, length int, size int) {
	c.table.AddVarchar(name, length, "")
	c.table.AddKey(name, size)
}

func (c *command) AddUniqueBigInt(name string) {
	c.table.AddBigInt(name, 20, 0, false)
	c.table.AddUniqueKey(name)
}

func (c *command) AddUniqueInt(name string) {
	c.table.AddInt(name, 11, 0, false)
	c.table.AddUniqueKey(name)
}

func (c *command) AddUniqueString(name string, length int) {
	c.table.AddVarchar(name, length, "")
	c.table.AddUniqueKey(name)
}

func (c *command) AddInt(name string) {
	c.table.AddInt(name, 11, 0, false)
}

func (c *command) AddBigInt(name string) {
	c.table.AddBigInt(name, 11, 0, false)
}

func (c *command) AddString(name string) {
	c.table.AddVarchar(name, 100, "")
}

func (c *command) AddStringWithLength(name string, length int) {
	c.table.AddVarchar(name, length, "")
}

func (c *command) AddText(name string) {
	c.table.AddText(name)
}

func (c *command) AddTimestamp(name string, onDefault, onUpdate bool) {
	c.table.AddTimestamp(name, onDefault, onUpdate)
}
