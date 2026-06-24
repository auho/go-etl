package buildtable

import (
	"errors"
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
	"github.com/auho/go-etl/v3/insight/assistant/sqlbuilder/ddl/command/mysql"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ Tabler = (*table)(nil)

type Tabler interface {
	GetTableName() string
	GetCommand() *schema.Command
	SQL() string
	Build() error
	ExecCommand(func(*schema.Command))

	withConfig(Config)
}

type table struct {
	*schema.Command
	config Config
	db     *simpledb.SimpleDB
}

func (t *table) initCommand(name string) {
	t.Command = &schema.Command{Table: mysql.NewTableSimple(name)}
}

func (t *table) GetCommand() *schema.Command {
	return t.Command
}

func (t *table) GetTableName() string {
	return t.Command.TableName()
}

func (t *table) SQL() string {
	return t.Command.SQLForCreate()
}

func (t *table) Build() error {
	if t.config.Recreate {
		err := t.db.Drop(t.TableName())
		if err != nil {
			return t.formatError(fmt.Errorf("drop error; %w", err))
		}
	} else if t.config.Truncate {
		err := t.db.Truncate(t.TableName())
		if err != nil {
			return t.formatError(fmt.Errorf("truncate error; %w", err))
		}
	}

	sql := t.SQL()
	if sql == "" {
		return t.formatError(errors.New("sql empty error"))
	}

	if t.db == nil {
		return t.formatError(errors.New("db empty error"))
	}

	err := t.db.GormDB().Exec(sql).Error
	if err != nil {
		return t.formatError(err)
	}

	return nil
}

func (t *table) ExecCommand(fn func(*schema.Command)) {
	fn(t.Command)
}

func (t *table) withConfig(config Config) {
	t.config = config
}

func (t *table) options(opts []TableOption) {
	for _, opt := range opts {
		opt(t)
	}
}

// exec model command
func (t *table) execRawCommandFunc(r assistant.Raw) {
	r.ExecCommand(t.Command)
}

func (t *table) formatError(err error) error {
	return fmt.Errorf("%s; %w", t.TableName(), err)
}
