package model

import (
	"fmt"
	"maps"
	"sort"

	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml"
	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml/command"
	"github.com/auho/go-etl/v2/job/means/tag"
)

var _ tag.Ruler = (*RuleItems)(nil)

// RuleItemsConfig
// rule items config
type RuleItemsConfig struct {
	Alias             map[string]string   // map[data name]output name
	Fixed             map[string]any      // map[key]value
	KeywordFormatFunc func(string) string // []func(data keyword value)regexp keyword value
}

func WithRuleItemsConfig(config RuleItemsConfig) func(*RuleItems) {
	return func(ri *RuleItems) {
		ri.alias = config.Alias
		ri.fixed = config.Fixed
		ri.keywordFormatFunc = config.KeywordFormatFunc
	}
}

// RuleItems
// alias: [data name] => [output name]
// fixed: [key] => [value]
// keywordFormatFunc: [data keyword value] => [regexp keyword value]
type RuleItems struct {
	rule              assistant.Ruler
	alias             map[string]string
	fixed             map[string]any
	keywordFormatFunc func(string) string
}

func NewRuleItems(rule assistant.Ruler, opts ...func(items *RuleItems)) *RuleItems {
	ri := &RuleItems{}
	ri.rule = rule

	for _, opt := range opts {
		opt(ri)
	}

	if ri.keywordFormatFunc == nil {
		ri.keywordFormatFunc = func(s string) string {
			return s
		}
	}

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
	_rule := ri.rule
	_originRule := ri.rule.ToOriginRule()

	selects := make(map[string]string)
	selects[_originRule.GetName()] = _rule.GetName()

	for _originLabel, _label := range _rule.LabelsAlias() {
		selects[_originLabel] = _label
	}

	selects[_originRule.KeywordName()] = _rule.KeywordName()

	table := dml.NewTable(_rule.TableName())

	for _originLabel, _label := range selects {
		if _labelAlias, _ok := ri.getAlias(_label); _ok {
			_label = _labelAlias
		}

		if _originLabel == _label {
			table = table.Select([]string{_originLabel})
		} else {
			table = table.SelectAlias(map[string]string{_originLabel: _label})
		}
	}

	table.OrderBy(_originRule.KeywordLenName(), command.SortDesc, _originRule.GetIdName(), command.SortAsc)

	var rows []map[string]any
	sql := table.Sql()
	err := _rule.GetDB().Raw(sql).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("rows error; %w", err)
	}

	var _newRows []map[string]string
	for _, row := range rows {
		_nRow := make(map[string]string, len(row))
		for _k, _v := range row {
			_nRow[_k] = _v.(string)
		}

		_newRows = append(_newRows, _nRow)
	}

	return _newRows, nil
}

func (ri *RuleItems) ItemsForRegexp() ([]map[string]string, error) {
	rows, err := ri.ItemsAlias()
	if err != nil {
		return nil, fmt.Errorf("ItemsAlias error; %w", err)
	}
	for i := range rows {
		rows[i][ri.KeywordNameAlias()] = ri.keywordFormatFunc(rows[i][ri.KeywordNameAlias()])
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

func (ri *RuleItems) Fixed() map[string]any {
	return ri.fixed
}

func (ri *RuleItems) FixedAlias() map[string]any {
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
