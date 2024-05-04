package model

import (
	"fmt"
	"sort"

	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ assistant.RuleConfigure = (*RuleConfig)(nil)

type RuleConfig struct {
	allowKeywordDuplicate bool
}

func (rc RuleConfig) AllowKeywordDuplicate() bool {
	return rc.allowKeywordDuplicate
}

type baseRule struct {
	model
	config        RuleConfig
	name          string // origin name
	length        int    // origin name length
	keywordLength int
	labels        map[string]int // map[label]label length

	alias             map[string]string // map[origin name][alias name]
	nameAlias         string            // alias of name
	labelsAlias       map[string]string // map[label]label alias
	labelsAliasLength map[string]int    // map[label alias]label alias length, alias labels

	independentTableName string // independent table name
}

func newBaseRule(name string, length, keywordLength int, labels map[string]int, db *simpleDb.SimpleDB) baseRule {
	br := baseRule{}
	br.name = name
	br.length = length
	br.keywordLength = keywordLength
	br.labels = labels
	br.db = db

	br.handlerAlias(nil)

	if br.length <= 0 {
		br.length = defaultStringLen
	}

	if br.keywordLength <= 0 {
		br.keywordLength = defaultStringLen
	}

	return br
}

func (br *baseRule) GetDB() *simpleDb.SimpleDB {
	return br.db
}

func (br *baseRule) GetName() string {
	return br.nameAlias
}

func (br *baseRule) GetNameLength() int {
	return br.length
}

func (br *baseRule) GetIdName() string {
	return "id"
}

func (br *baseRule) GetKeywordLength() int {
	return br.keywordLength
}

func (br *baseRule) GetLabels() map[string]int {
	return br.labelsAliasLength
}

func (br *baseRule) TagsName() []string {
	var tagsName []string
	tagsName = append(tagsName, br.GetName())
	tagsName = append(tagsName, br.LabelsName()...)

	return tagsName
}

func (br *baseRule) LabelsName() []string {
	var labels []string
	for label, _ := range br.GetLabels() {
		labels = append(labels, label)
	}

	sort.SliceIsSorted(labels, func(i, j int) bool {
		return labels[i] < labels[j]
	})

	return labels
}

func (br *baseRule) LabelsAlias() map[string]string {
	return br.labelsAlias
}

func (br *baseRule) LabelNumName() string {
	return fmt.Sprintf("%s_%s", br.nameAlias, NameLabelNum)
}

func (br *baseRule) KeywordName() string {
	return fmt.Sprintf("%s_%s", br.nameAlias, NameKeyword)
}

func (br *baseRule) KeywordLenName() string {
	return fmt.Sprintf("%s_%s", br.nameAlias, NameKeywordLen)
}

func (br *baseRule) KeywordNumName() string {
	return fmt.Sprintf("%s_%s", br.nameAlias, NameKeywordNum)
}

func (br *baseRule) KeywordAmountName() string {
	return fmt.Sprintf("%s_%s", br.nameAlias, NameKeywordAmount)
}

func (br *baseRule) Config() assistant.RuleConfigure {
	return br.config
}

func (br *baseRule) WithCommand(fn func(command *tablestructure.Command)) *baseRule {
	br.withCommand(fn)

	return br
}

func (br *baseRule) handlerAlias(alias map[string]string) {
	br.alias = alias

	if v, ok := alias[br.name]; ok {
		br.nameAlias = v
	} else {
		br.nameAlias = br.name
	}

	_labelsAlias := make(map[string]string, len(br.labels))
	_labelsAliasLength := make(map[string]int, len(br.labels))
	for label, length := range br.labels {
		if v, ok := alias[label]; ok {
			_labelsAlias[label] = v
			_labelsAliasLength[v] = length
		} else {
			_labelsAlias[label] = label
			_labelsAliasLength[label] = length
		}
	}

	br.labelsAlias = _labelsAlias
	br.labelsAliasLength = _labelsAliasLength
}
