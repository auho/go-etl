package model

import (
	"fmt"

	"github.com/auho/go-etl/tool"
)

type Rule struct {
	name string
}

func NewRule(name string) *Rule {
	r := &Rule{}
	r.name = name

	return r
}

func (r *Rule) GetName() string {
	return r.name
}

func (r *Rule) TableName() string {
	return fmt.Sprintf("%s_%s", tool.RuleTableNamePrefix, r.name)
}

func (r *Rule) Keyword() string {
	return fmt.Sprintf("%s_keyword", r.name)
}

func (r *Rule) KeywordNum() string {
	return fmt.Sprintf("%s_keyword_num", r.name)
}

func (r *Rule) DataTableName(n string) string {
	return fmt.Sprintf("%s_%s_%s", tool.RuleTableNamePrefix, n, r.name)
}
