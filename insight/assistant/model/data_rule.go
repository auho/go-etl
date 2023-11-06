package model

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant"
)

var _ assistant.Ruler = (*Rule)(nil)

type DataRule struct {
	*Rule
	data *Data
}

func NewDataRule(data *Data, rule *Rule) *DataRule {
	dr := &DataRule{}
	dr.data = data
	dr.Rule = rule

	return dr
}

func (dr *DataRule) GetData() *Data {
	return dr.data
}

func (dr *DataRule) GetRule() *Rule {
	return dr.Rule
}

func (dr *DataRule) TableName() string {
	return fmt.Sprintf("%s_%s_%s", NameRule, dr.data.GetName(), dr.GetName())
}
