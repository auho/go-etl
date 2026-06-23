package schema

import (
	"github.com/auho/go-etl/v2/insight/assistant/sqlbuilder/ddl/command/mysql"
)

type Command struct {
	Table *mysql.Table
}

func NewCommandMysql() *Command {
	return &Command{Table: mysql.NewTable()}
}

func (c *Command) TableName() string {
	return c.Table.GetName()
}

func (c *Command) SQLForCreate() string {
	return c.Table.SQLForCreate()
}

func (c *Command) SQLForAlterAdd() []string {
	return c.Table.SQLForAlterAdd()
}

func (c *Command) SQLForAlterChange() []string {
	return c.Table.SQLForAlterChange()
}

func (c *Command) AddKey(name string, size int) {
	c.Table.AddKey(name, size)
}

func (c *Command) AddPK(name string) {
	c.Table.AddPK(name)
}

func (c *Command) AddPKBigInt(name string) {
	c.Table.AddPKBigInt(name, 20)
}

func (c *Command) AddPKInt(name string) {
	c.Table.AddPKInt(name, 11)
}

func (c *Command) AddPKString(name string, length int) {
	c.Table.AddPKVarchar(name, length)
}

func (c *Command) AddKeyBigInt(name string) {
	c.Table.AddBigInt(name, 20, 0, false)
	c.Table.AddKey(name, 0)
}

func (c *Command) AddKeyInt(name string) {
	c.Table.AddInt(name, 11, 0, false)
	c.Table.AddKey(name, 0)
}

func (c *Command) AddKeyString(name string, length int, size int) {
	c.Table.AddVarchar(name, length, "")
	c.Table.AddKey(name, size)
}

func (c *Command) AddUniqueBigInt(name string) {
	c.Table.AddBigInt(name, 20, 0, false)
	c.Table.AddUniqueKey(name)
}

func (c *Command) AddUniqueInt(name string) {
	c.Table.AddInt(name, 11, 0, false)
	c.Table.AddUniqueKey(name)
}

func (c *Command) AddUniqueString(name string, length int) *mysql.Field {
	f := c.Table.AddVarchar(name, length, "")
	c.Table.AddUniqueKey(name)

	return f
}

func (c *Command) AddInt(name string) {
	c.Table.AddInt(name, 11, 0, false)
}

func (c *Command) AddBigInt(name string) {
	c.Table.AddBigInt(name, 11, 0, false)
}

func (c *Command) AddString(name string) {
	c.Table.AddVarchar(name, 30, "")
}

func (c *Command) AddStringWithLength(name string, length int) *mysql.Field {
	return c.Table.AddVarchar(name, length, "")
}

func (c *Command) AddText(name string) {
	c.Table.AddText(name)
}

func (c *Command) AddTimestamp(name string, onDefault, onUpdate bool) {
	c.Table.AddTimestamp(name, onDefault, onUpdate)
}

func (c *Command) AddDecimal(name string, m, d int) {
	c.Table.AddDecimal(name, m, d, 0)
}
