package altertable

import (
	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpledb "github.com/auho/go-simple-db/v2"
)

var _db *simpledb.SimpleDB
var _raw assistant.Rawer

func ExampleNewTable() {
	_ = NewTable("tableName").
		WithCommand(func(command *tablestructure.Command) {
			command.AddInt("int1")
			command.AddString("s1")
		}).Build(_db)

	_t := NewTable("tableName")
	_t.AddInt("int2")
	_t.AddString("s2")
	_ = _t.BuildChange(_db)
}

func ExampleNewModelTable() {
	_ = NewModelTable(_raw).
		WithCommand(func(command *tablestructure.Command) {
			command.AddInt("int1")
		}).Build()

	_t := NewModelTable(_raw)
	_t.AddInt("int2")
	_t.AddString("s2")
	_ = _t.BuildChange()
}
