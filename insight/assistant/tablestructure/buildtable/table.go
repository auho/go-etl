package buildtable

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/accessory/ddl/command/mysql"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ Tabler = (*table)(nil)

type Tabler interface {
	GetTableName() string
	GetCommand() *tablestructure.Command
	Sql() string
	Build() error

	withCommand(func(*tablestructure.Command))
	withConfig(Config)
}

type Config struct {
	Recreate bool
	Truncate bool
}

type table struct {
	*tablestructure.Command
	config     Config
	db         *simpleDb.SimpleDB
	commandFun func(*tablestructure.Command)
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

func (t *table) Build() error {
	if t.config.Recreate {
		err := t.db.Drop(t.TableName())
		if err != nil {
			return fmt.Errorf("drop error; %w", err)
		}
	} else if t.config.Truncate {
		err := t.db.Truncate(t.TableName())
		if err != nil {
			return fmt.Errorf("truncate error; %w", err)
		}
	}

	sql := t.Sql()
	if sql == "" {
		return fmt.Errorf("sql empty error")
	}

	if t.db == nil {
		return fmt.Errorf("db empty error")
	}

	return t.db.Exec(sql).Error
}

func (t *table) withConfig(config Config) {
	t.config = config
}

func (t *table) withCommand(fn func(command *tablestructure.Command)) {
	t.commandFun = fn
}

func (t *table) options(opts []TableOption) {
	for _, opt := range opts {
		opt(t)
	}
}

func (t *table) execCommand() {
	if t.commandFun != nil {
		t.commandFun(t.Command)
	}
}

func (t *table) execRowsCommand(r assistant.Rowsor) {
	r.ExecCommand(t.Command)
}
