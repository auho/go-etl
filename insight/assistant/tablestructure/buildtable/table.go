package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant/accessory/ddl/command/mysql"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ Tabler = (*table)(nil)

type Tabler interface {
	GetTableName() string
	Sql() string
	Build(db *simpleDb.SimpleDB) error
}

type table struct {
	*tablestructure.Command
}

func (t *table) initCommand(name string) {
	t.Command = &tablestructure.Command{Table: &mysql.Table{}}
	t.Command.Table.SetName(name).SetEngineMyISAM()
}

func (t *table) GetCommand() *tablestructure.Command {
	return t.Command
}

func (t *table) GetTableName() string {
	return t.Command.TableName()
}

func (t *table) Sql() string {
	return t.Command.SqlForCreate()
}

func (t *table) Build(db *simpleDb.SimpleDB) error {
	sql := t.Sql()

	return db.Exec(sql).Error
}
