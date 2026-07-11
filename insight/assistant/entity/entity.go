package entity

import (
	"github.com/auho/go-etl/v3/insight/assistant/schema"
	simpledb "github.com/auho/go-simple-db/v3"
)

type base struct {
	commandFunc func(command *schema.Command)
	db          *simpledb.SimpleDB
}

func (m *base) withCommand(fn func(command *schema.Command)) {
	m.commandFunc = fn
}

// ExecCommand
// exec base table command
func (m *base) ExecCommand(command *schema.Command) {
	if m.commandFunc != nil {
		m.commandFunc(command)
	}
}
