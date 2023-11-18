package model

import (
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
)

func ExampleNewRule() {
	_ = NewRule("one", 20, 20, map[string]int{"label1": 20, "label2": 30}, nil).
		WithCommand(func(command *tablestructure.Command) {
			command.AddString("field1")
		})
}

func ExampleNewRuleSimple() {
	_ = NewRuleSimple("one", []string{"label1", "label2"}, nil).
		WithCommand(func(command *tablestructure.Command) {
			command.AddString("field1")
		})
}

func ExampleNewRuleItems() {
	_rule := NewRuleSimple("one", []string{"label1", "label2"}, nil)

	_ = NewRuleItems(
		_rule,
		WithRuleItemsConfig(RuleItemsConfig{
			Alias:             map[string]string{"label1": "label1_alias", "label2": "label2_alias"},
			Fixed:             map[string]any{"fixed1": 1},
			KeywordFormatFunc: func(s string) string { return s },
		}),
	)
}
