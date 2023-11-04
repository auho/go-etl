package table

import (
	"github.com/auho/go-etl/v2/insight/model/ddl/command/mysql"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ Tabler = (*table)(nil)

type Tabler interface {
	GetCommand() *command
	GetTableName() string
	Sql() string
	Build(db *simpleDb.SimpleDB) error
}

type table struct {
	*command
}

func (t *table) initCommand(name string) {
	t.command = &command{table: &mysql.Table{}}
	t.command.table.SetName(name).SetEngineMyISAM()
}

func (t *table) GetCommand() *command {
	return t.command
}

func (t *table) GetTableName() string {
	return t.command.TableName()
}

func (t *table) Sql() string {
	return t.command.Sql()
}

func (t *table) Build(db *simpleDb.SimpleDB) error {
	sql := t.Sql()

	return db.Exec(sql).Error
}
