package model

import (
	"fmt"

	go_etl "github.com/auho/go-etl"
)

type Rule struct {
	name string
}

func NewRule(name string) *Rule {
	r := &Rule{}
	r.name = name

	return r
}

func (r *Rule) getName() string {
	return r.name
}

func (r *Rule) tableName() string {
	return fmt.Sprintf("%s_%s", go_etl.RuleTableNamePrefix, r.name)
}

func (r *Rule) keyword() string {
	return fmt.Sprintf("%s_keyword", r.name)
}

func (r *Rule) keywordNum() string {
	return fmt.Sprintf("%s_keyword_num", r.name)
}
