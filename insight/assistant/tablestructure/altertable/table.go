package altertable

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant/accessory/ddl/command/mysql"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpleDb "github.com/auho/go-simple-db/v2"
)

type Table struct {
	*tablestructure.Command
}

func NewTable(tableName string) *Table {
	t := &Table{}
	t.Command = &tablestructure.Command{Table: &mysql.Table{}}
	t.Command.Table.SetName(tableName)

	return t
}

func (t *Table) GetTableName() string {
	return t.Command.TableName()
}

func (t *Table) Sql() []string {
	return t.Command.SqlForAlterAdd()
}

func (t *Table) Exec(fn func(command *tablestructure.Command)) *Table {
	fn(t.Command)

	return t
}

func (t *Table) Build(db *simpleDb.SimpleDB) error {
	for _, sql := range t.Sql() {
		err := db.Exec(sql).Error
		if err != nil {
			return fmt.Errorf("exec error; %w", err)
		}
	}

	return nil
}
