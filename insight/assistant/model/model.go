package model

import (
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpledb "github.com/auho/go-simple-db/v2"
)

type model struct {
	commandFunc func(command *tablestructure.Command)
	db          *simpledb.SimpleDB
}

func (m *model) withCommand(fn func(command *tablestructure.Command)) {
	m.commandFunc = fn
}

// ExecCommand
// exec model table command
func (m *model) ExecCommand(command *tablestructure.Command) {
	if m.commandFunc != nil {
		m.commandFunc(command)
	}
}
