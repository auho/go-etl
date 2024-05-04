package altertable

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpleDb "github.com/auho/go-simple-db/v2"
)

type baseTable struct {
	*tablestructure.Command
	commandFun func(*tablestructure.Command)
}

func newBaseTable(tableName string) baseTable {
	bt := baseTable{
		Command: nil,
	}

	bt.Command = tablestructure.NewCommandMysql()
	bt.Command.Table.SetName(tableName)

	return bt
}

func (bt *baseTable) GetTableName() string {
	return bt.Command.TableName()
}

func (bt *baseTable) Sql() []string {
	bt.execCommand()
	return bt.Command.SqlForAlterAdd()
}

func (bt *baseTable) SqlForChange() []string {
	bt.execCommand()
	return bt.Command.SqlForAlterChange()
}

func (bt *baseTable) build(sqls []string, db *simpleDb.SimpleDB) error {
	for _, sql := range sqls {
		err := db.Exec(sql).Error
		if err != nil {
			return fmt.Errorf("build[%s] exec error; %w", bt.TableName(), err)
		}
	}

	return nil
}

func (bt *baseTable) execCommand() {
	if bt.commandFun != nil {
		bt.commandFun(bt.Command)
	}
}
