package alter

import (
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant/schema"
	simpledb "github.com/auho/go-simple-db/v3"
)

type baseTable struct {
	*schema.Command
	commandFunc func(*schema.Command)
}

func newBaseTable(tableName string) baseTable {
	bt := baseTable{
		Command: nil,
	}

	bt.Command = schema.NewCommandMysql()
	bt.Command.Table.SetName(tableName)

	return bt
}

func (bt *baseTable) GetTableName() string {
	return bt.Command.TableName()
}

func (bt *baseTable) SQL() []string {
	bt.execCommand()
	return bt.Command.SQLForAlterAdd()
}

func (bt *baseTable) SqlForChange() []string {
	bt.execCommand()
	return bt.Command.SQLForAlterChange()
}

func (bt *baseTable) build(sqls []string, db *simpledb.SimpleDB) error {
	for _, sql := range sqls {
		err := db.GormDB().Exec(sql).Error
		if err != nil {
			return fmt.Errorf("exec[%s]: %w", bt.TableName(), err)
		}
	}

	return nil
}

func (bt *baseTable) execCommand() {
	if bt.commandFunc != nil {
		bt.commandFunc(bt.Command)
	}
}
