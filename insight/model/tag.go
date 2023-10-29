package model

import (
	"fmt"
)

type Tag struct {
	data *Data
	rule *Rule
}

func NewTag(data *Data, rule *Rule) *Tag {
	t := &Tag{}
	t.data = data
	t.rule = rule

	return t
}

func (t *Tag) GetData() *Data {
	return t.data
}

func (t *Tag) GetRule() *Rule {
	return t.rule
}

func (t *Tag) TableName() string {
	return fmt.Sprintf("%s_%s_%s", NameTag, t.data.GetName(), t.rule.GetName())
}

func (t *Tag) KeywordNum() string {
	return fmt.Sprintf("%s_%s", t.rule.GetName(), NameKeywordNum)
}
