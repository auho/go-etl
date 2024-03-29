package altertable

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant/accessory/ddl/command/mysql"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpleDb "github.com/auho/go-simple-db/v2"
)

type Table struct {
	*tablestructure.Command
	commandFun func(*tablestructure.Command)
}

func NewTable(tableName string) *Table {
	t := &Table{}
	t.Command = &tablestructure.Command{Table: mysql.NewTable()}
	t.Command.Table.SetName(tableName)

	return t
}

func (t *Table) execCommand() {
	if t.commandFun != nil {
		t.commandFun(t.Command)
	}
}

func (t *Table) GetTableName() string {
	return t.Command.TableName()
}

func (t *Table) Sql() []string {
	t.execCommand()
	return t.Command.SqlForAlterAdd()
}

func (t *Table) SqlForChange() []string {
	t.execCommand()
	return t.Command.SqlForAlterChange()
}

func (t *Table) WithCommand(fn func(command *tablestructure.Command)) *Table {
	t.commandFun = fn

	return t
}

func (t *Table) Build(db *simpleDb.SimpleDB) error {
	return t.build(t.Sql(), db)
}

func (t *Table) BuildChange(db *simpleDb.SimpleDB) error {
	return t.build(t.SqlForAlterChange(), db)
}

func (t *Table) build(sqls []string, db *simpleDb.SimpleDB) error {
	for _, sql := range sqls {
		err := db.Exec(sql).Error
		if err != nil {
			return fmt.Errorf("build exec error; %w", err)
		}
	}

	return nil
}
