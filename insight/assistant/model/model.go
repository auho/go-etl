package model

import (
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpleDb "github.com/auho/go-simple-db/v2"
)

type model struct {
	commandFun func(command *tablestructure.Command)
	db         *simpleDb.SimpleDB
}

func (m *model) withCommand(fn func(command *tablestructure.Command)) {
	m.commandFun = fn
}

func (m *model) execCommand(command *tablestructure.Command) {
	if m.commandFun != nil {
		m.commandFun(command)
	}
}
