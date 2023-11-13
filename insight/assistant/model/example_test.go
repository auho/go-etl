package model

func ExampleNewRuleItems() {
	_rule := NewRuleSimple("one", []string{"label1", "label2"}, nil)
	_ = NewRuleItems(_rule, WithRuleItemsConfig(RuleItemsConfig{
		Alias:             map[string]string{"label1": "label1_alias", "label2": "label2_alias"},
		Fixed:             map[string]any{"fixed1": 1},
		KeywordFormatFunc: func(s string) string { return s },
	}))
}
