package model

import (
	"fmt"
	"maps"
	"sort"

	"github.com/auho/go-etl/v2/insight/model/dml"
	"github.com/auho/go-etl/v2/insight/model/dml/command"
	"github.com/auho/go-etl/v2/means/tag"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ tag.Ruler = (*RuleItems)(nil)

type RuleItems struct {
	rule  Ruler
	alias map[string]string
	fixed map[string]string

	db *simpleDb.SimpleDB
}

func NewRuleItems(db *simpleDb.SimpleDB, rule Ruler, alias, fixed map[string]string) *RuleItems {
	ri := &RuleItems{}
	ri.db = db
	ri.rule = rule
	ri.alias = alias
	ri.fixed = fixed

	return ri
}

func (ri *RuleItems) getAlias(s string) (string, bool) {
	if ns, ok := ri.alias[s]; ok {
		return ns, ok
	} else {
		return s, ok
	}
}

func (ri *RuleItems) Items() ([]map[string]string, error) {
	fields := []string{ri.rule.GetName(), ri.rule.KeywordName()}
	fields = append(fields, ri.rule.LabelsName()...)

	table := dml.NewTable(ri.TableName())
	for label, _ := range ri.rule.GetLabels() {
		if labelAlias, ok := ri.getAlias(label); ok {
			table = table.SelectAlias(map[string]string{label: labelAlias})
		} else {
			table = table.Select([]string{label})
		}
	}

	table.OrderBy(map[string]string{
		ri.rule.KeywordLenName(): command.SortDesc,
		ri.rule.GetIdName():      command.SortASC,
	})

	var rows []map[string]string
	err := ri.db.Raw(table.Sql()).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("rows error; %w", err)
	}

	return rows, nil
}

func (ri *RuleItems) Name() string {
	return ri.rule.GetName()
}

func (ri *RuleItems) TableName() string {
	return ri.rule.TableName()
}

func (ri *RuleItems) Labels() []string {
	var labels []string
	for label, _ := range ri.rule.GetLabels() {
		labels = append(labels, label)
	}

	sort.Slice(labels, func(i, j int) bool {
		return labels[i] < labels[j]
	})

	return labels
}

func (ri *RuleItems) KeywordName() string {
	return ri.rule.KeywordName()
}

func (ri *RuleItems) KeywordLenName() string {
	return ri.rule.KeywordLenName()
}

func (ri *RuleItems) KeywordNumName() string {
	return ri.rule.KeywordNumName()
}

func (ri *RuleItems) Fixed() map[string]string {
	fixed := maps.Clone(ri.fixed)

	for k, v := range fixed {
		if nk, ok := ri.getAlias(k); ok {
			fixed[nk] = v
			delete(fixed, k)
		}
	}

	return fixed
}

func (ri *RuleItems) FixedKeys() []string {
	var keys []string
	for k := range ri.fixed {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	return keys
}
