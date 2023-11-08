package model

import (
	"fmt"
	"maps"
	"sort"

	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ assistant.Ruler = (*Rule)(nil)

const defaultStringLen = 30

type Rule struct {
	model
	name          string // origin name
	length        int    // origin name length
	keywordLength int
	labels        map[string]int

	nameAlias   string         // alias name
	labelsAlias map[string]int // alias labels
}

func NewRuleSimple(name string, labels []string, db *simpleDb.SimpleDB) *Rule {
	_labels := make(map[string]int, len(labels))
	for _, label := range labels {
		_labels[label] = defaultStringLen
	}

	return NewRule(name, defaultStringLen, defaultStringLen, _labels, db)
}

func NewRule(name string, length, keywordLength int, labels map[string]int, db *simpleDb.SimpleDB) *Rule {
	r := &Rule{}
	r.name = name
	r.length = length
	r.keywordLength = keywordLength
	r.labels = labels
	r.db = db

	r.nameAlias = r.name
	r.labelsAlias = maps.Clone(r.labels)

	if r.length <= 0 {
		r.length = defaultStringLen
	}

	if r.keywordLength <= 0 {
		r.keywordLength = defaultStringLen
	}

	return r
}

func (r *Rule) GetDB() *simpleDb.SimpleDB {
	return r.db
}

func (r *Rule) GetName() string {
	return r.nameAlias
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
	return r.labelsAlias
}

func (r *Rule) LabelsName() []string {
	var labels []string
	for label, _ := range r.GetLabels() {
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
	return fmt.Sprintf("%s_%s", r.nameAlias, NameKeyword)
}

func (r *Rule) KeywordLenName() string {
	return fmt.Sprintf("%s_%s", r.nameAlias, NameKeywordLen)
}

func (r *Rule) KeywordNumName() string {
	return fmt.Sprintf("%s_%s", r.nameAlias, NameKeywordNum)
}

func (r *Rule) WithCommand(fn func(command *tablestructure.Command)) *Rule {
	r.withCommand(fn)

	return r
}

func (r *Rule) ToOriginRule() assistant.Ruler {
	return NewRule(r.name, r.length, r.keywordLength, maps.Clone(r.labels), r.db)
}

func (r *Rule) ToAliasRule(alias map[string]string) *Rule {
	_rule := NewRule(r.name, r.length, r.keywordLength, maps.Clone(r.labels), r.db)

	if v, ok := alias[_rule.name]; ok {
		_rule.nameAlias = v
	} else {
		_rule.nameAlias = _rule.name
	}

	_labels := make(map[string]int, len(_rule.labels))
	for label, length := range _rule.labels {
		if v, ok := alias[label]; ok {
			_labels[v] = length
		} else {
			_labels[label] = length
		}
	}

	_rule.labelsAlias = _labels

	return _rule
}

func (r *Rule) ToItems(opts ...func(items *RuleItems)) *RuleItems {
	return NewRuleItems(r, opts...)
}
