package model

import (
	"fmt"
	"sort"

	"github.com/auho/go-etl/v2/insight/assistant"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ assistant.Ruler = (*Rule)(nil)

type Rule struct {
	name          string
	length        int
	keywordLength int
	labels        map[string]int
	db            *simpleDb.SimpleDB
}

func NewRuleSimple(name string, labels map[string]int, db *simpleDb.SimpleDB) *Rule {
	return NewRule(name, 30, 30, labels, db)
}

func NewRule(name string, length, keywordLength int, labels map[string]int, db *simpleDb.SimpleDB) *Rule {
	r := &Rule{}
	r.name = name
	r.length = length
	r.keywordLength = keywordLength
	r.labels = labels
	r.db = db

	if r.length <= 0 {
		r.length = 30
	}

	if r.keywordLength <= 0 {
		r.keywordLength = 30
	}

	return r
}

func (r *Rule) GetDB() *simpleDb.SimpleDB {
	return r.db
}

func (r *Rule) GetName() string {
	return r.name
}

func (r *Rule) GetNameLength() int {
	return r.length
}

func (r *Rule) GetIdName() string {
	return "id"
}

func (r *Rule) GetKeywordLength() int {
	return r.keywordLength
}

func (r *Rule) GetLabels() map[string]int {
	return r.labels
}

func (r *Rule) LabelsName() []string {
	var labels []string
	for label, _ := range r.labels {
		labels = append(labels, label)
	}

	sort.SliceIsSorted(labels, func(i, j int) bool {
		return labels[i] < labels[j]
	})

	return labels
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

func (r *Rule) KeywordNumName() string {
	return fmt.Sprintf("%s_%s", r.name, NameKeywordNum)
}
