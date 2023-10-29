package model

import (
	"fmt"
)

type Rule struct {
	name          string
	length        int
	keywordLength int
	labels        map[string]int
}

func NewRule(name string, length, keywordLength int, labels map[string]int) *Rule {
	r := &Rule{}
	r.name = name
	r.length = length
	r.keywordLength = keywordLength
	r.labels = labels

	if r.length <= 0 {
		r.length = 30
	}

	if r.keywordLength <= 0 {
		r.keywordLength = 30
	}

	return r
}

func (r *Rule) GetName() string {
	return r.name
}

func (r *Rule) GetLength() int {
	return r.length
}

func (r *Rule) GetKeywordLength() int {
	return r.keywordLength
}

func (r *Rule) GetLabels() map[string]int {
	return r.labels
}

func (r *Rule) TableName() string {
	return fmt.Sprintf("%s_%s", NameRule, r.name)
}

func (r *Rule) KeywordName() string {
	return fmt.Sprintf("%s_%s", r.name, NameKeyword)
}

func (r *Rule) KeywordLenName() string {
	return fmt.Sprintf("%s_%s", r.name, NameKeywordLen)
}
