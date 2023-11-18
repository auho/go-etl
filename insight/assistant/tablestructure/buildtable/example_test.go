package buildtable

import (
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
)

func ExampleNewTagDataRuleTable() {
	_ = NewTagDataRuleTable(
		nil,
		WithConfig(Config{
			Recreate: true,
			Truncate: true,
		})).
		WithCommand(func(command *tablestructure.Command) {
			// add columns
			command.AddString("a")
			command.AddString("b")
		}).Build()
}

func ExampleNewTagDataRulesTable() {
	_ = NewTagDataRulesTable(
		nil,
		WithConfig(Config{
			Recreate: true,
			Truncate: true,
		})).
		WithCommand(func(command *tablestructure.Command) {
			// add columns
			command.AddString("a")
			command.AddString("b")
		}).
		Build()
}
