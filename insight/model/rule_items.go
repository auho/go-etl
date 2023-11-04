package model

import (
	"fmt"
	"maps"
	"sort"

	"github.com/auho/go-etl/v2/insight/model/dml"
	"github.com/auho/go-etl/v2/insight/model/dml/command"
	"github.com/auho/go-etl/v2/means/tag"
)

var _ tag.Ruler = (*RuleItems)(nil)

type RuleItems struct {
	rule              Ruler
	alias             map[string]string
	fixed             map[string]string
	keywordFormatFunc []func(string) string
}

func NewRuleItems(rule Ruler, alias, fixed map[string]string, keywordFormatFunc []func(string) string) *RuleItems {
	ri := &RuleItems{}
	ri.rule = rule
	ri.alias = alias
	ri.fixed = fixed
	ri.keywordFormatFunc = keywordFormatFunc

	return ri
}

func (ri *RuleItems) getAlias(s string) (string, bool) {
	if ns, ok := ri.alias[s]; ok {
		return ns, ok
	} else {
		return s, ok
	}
}

func (ri *RuleItems) ItemsAlias() ([]map[string]string, error) {
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
	err := ri.rule.GetDB().Raw(table.Sql()).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("rows error; %w", err)
	}

	return rows, nil
}

func (ri *RuleItems) ItemsForRegexp() ([]map[string]string, error) {
	rows, err := ri.ItemsAlias()
	if err != nil {
		return nil, fmt.Errorf("ItemsAlias error; %w", err)
	}
	for i := range rows {
		for _, f := range ri.keywordFormatFunc {
			rows[i][ri.KeywordNameAlias()] = f(rows[i][ri.KeywordNameAlias()])
		}
	}

	return rows, nil
}

func (ri *RuleItems) Name() string {
	return ri.rule.GetName()
}

func (ri *RuleItems) NameAlias() string {
	s, _ := ri.getAlias(ri.Name())
	return s
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

func (ri *RuleItems) LabelsAlias() []string {
	var labels []string
	for _, label := range ri.Labels() {
		label, _ = ri.getAlias(label)

		labels = append(labels, label)
	}

	return labels
}

func (ri *RuleItems) KeywordName() string {
	return ri.rule.KeywordName()
}

func (ri *RuleItems) KeywordNameAlias() string {
	s, _ := ri.getAlias(ri.KeywordName())

	return s
}

func (ri *RuleItems) KeywordLenName() string {
	return ri.rule.KeywordLenName()
}

func (ri *RuleItems) KeywordLenNameAlias() string {
	s, _ := ri.getAlias(ri.KeywordLenName())

	return s
}

func (ri *RuleItems) KeywordNumName() string {
	return ri.rule.KeywordNumName()
}

func (ri *RuleItems) KeywordNumNameAlias() string {
	s, _ := ri.getAlias(ri.KeywordNumName())

	return s
}

func (ri *RuleItems) Fixed() map[string]string {
	return ri.fixed
}

func (ri *RuleItems) FixedAlias() map[string]string {
	fixed := maps.Clone(ri.Fixed())

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
	for k := range ri.Fixed() {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	return keys
}

func (ri *RuleItems) FixedKeysAlias() []string {
	var keys []string
	for k := range ri.FixedAlias() {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	return keys
}
