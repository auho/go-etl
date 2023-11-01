package model

import (
	"fmt"
)

type Tag struct {
	data Dataor
	rule Ruler
}

func NewTag(data Dataor, rule Ruler) *Tag {
	t := &Tag{}
	t.data = data
	t.rule = rule

	return t
}

func (t *Tag) GetData() Dataor {
	return t.data
}

func (t *Tag) GetRule() Ruler {
	return t.rule
}

func (t *Tag) TableName() string {
	return fmt.Sprintf("%s_%s_%s", NameTag, t.data.GetName(), t.rule.GetName())
}

func (t *Tag) KeywordNumName() string {
	return fmt.Sprintf("%s_%s", t.rule.GetName(), NameKeywordNum)
}
