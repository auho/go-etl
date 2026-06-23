package model

import (
	"github.com/auho/go-etl/v2/insight/assistant/schema"
	simpledb "github.com/auho/go-simple-db/v2"
)

type model struct {
	commandFunc func(command *schema.Command)
	db          *simpledb.SimpleDB
}

func (m *model) withCommand(fn func(command *schema.Command)) {
	m.commandFunc = fn
}

// ExecCommand
// exec model table command
func (m *model) ExecCommand(command *schema.Command) {
	if m.commandFunc != nil {
		m.commandFunc(command)
	}
}
