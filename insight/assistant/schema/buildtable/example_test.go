package buildtable

import (
	"github.com/auho/go-etl/v3/insight/assistant/schema"
)

func ExampleNewTagDataRuleTable() {
	_ = NewTagDataRuleTable(
		nil,
		WithConfig(Config{
			Recreate: true,
			Truncate: true,
		})).
		WithCommand(func(command *schema.Command) {
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
		WithCommand(func(command *schema.Command) {
			// add columns
			command.AddString("a")
			command.AddString("b")
		}).
		Build()
}
