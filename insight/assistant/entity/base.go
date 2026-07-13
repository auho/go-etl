package entity

import (
	"github.com/auho/go-etl/v3/insight/assistant/schema"
	simpledb "github.com/auho/go-simple-db/v3"
)

// base is the foundational embedded struct for all entity types.
// It provides a shared database connection and a command hook mechanism
// that allows callers to customize database schema operations (e.g., adding columns).
type base struct {
	commandFunc func(command *schema.Command) // optional hook to modify schema commands before execution
	db          *simpledb.SimpleDB            // shared database connection
}

// DB returns the database connection.
func (b *base) DB() *simpledb.SimpleDB {
	return b.db
}

// withCommand sets the command hook that will be called before ExecCommand.
func (b *base) withCommand(fn func(command *schema.Command)) {
	b.commandFunc = fn
}

// ExecCommand executes a schema command, invoking the command hook if set.
// This allows customization of the table schema before it is applied.
func (b *base) ExecCommand(command *schema.Command) {
	if b.commandFunc != nil {
		b.commandFunc(command)
	}
}
