package model

import (
	"fmt"
)

type DataRule struct {
	data *Data
	rule *Rule
}

func NewDataRule(data *Data, rule *Rule) *DataRule {
	dr := &DataRule{}
	dr.data = data
	dr.rule = rule

	return dr
}

func (dr *DataRule) GetData() *Data {
	return dr.data
}

func (dr *DataRule) GetRule() *Rule {
	return dr.rule
}

func (dr *DataRule) TableName() string {
	return fmt.Sprintf("%s_%s_%s", NameRule, dr.data.GetName(), dr.rule.GetName())
}

func (dr *DataRule) KeywordName() string {
	return fmt.Sprintf("%s_%s", dr.rule.GetName(), NameKeyword)
}

func (dr *DataRule) KeywordLenName() string {
	return fmt.Sprintf("%s_%s", dr.rule.GetName(), NameKeywordLen)
}
