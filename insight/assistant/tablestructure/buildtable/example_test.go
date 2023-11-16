package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
)

func ExampleNewTagDataRuleTable() {
	_ = NewTagDataRuleTable(
		nil,
		WithCommand(func(command *tablestructure.Command) {
			// add columns
			command.AddString("a")
			command.AddString("b")
		}),
		WithConfig(Config{
			Recreate: true,
			Truncate: true,
		})).Build()
}

func ExampleNewTagDataRulesTable() {
	_ = NewTagDataRulesTable(
		nil,
		WithCommand(func(command *tablestructure.Command) {
			// add columns
			command.AddString("a")
			command.AddString("b")
		}),
		WithConfig(Config{
			Recreate: true,
			Truncate: true,
		})).Build()
}
