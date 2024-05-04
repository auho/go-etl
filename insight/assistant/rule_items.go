package assistant

import (
	"fmt"
	"sort"

	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml"
	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml/command"
	"github.com/auho/go-etl/v2/job/means"
)

var _ means.Ruler = (*RuleItems)(nil)

// RuleItemsConfig
// rule items config
type RuleItemsConfig struct {
	Alias             map[string]string   // map[data name]output name
	KeywordFormatFunc func(string) string // []func(data keyword value)regexp keyword value
}

func WithRuleItemsConfig(config RuleItemsConfig) func(*RuleItems) {
	return func(ri *RuleItems) {
		ri.alias = config.Alias
		ri.keywordFormatFunc = config.KeywordFormatFunc
	}
}

// RuleItems
// alias: [data name] => [output name]
// keywordFormatFunc: [data keyword value] => [regexp keyword value]
type RuleItems struct {
	rule              Ruler
	alias             map[string]string
	keywordFormatFunc func(string) string

	tableName string
	name      string
	nameAlias string

	tags        []string
	labels      []string
	tagsAlias   []string
	labelsAlias []string

	labelNumName           string
	labelNumNameAlias      string
	keywordName            string
	keywordNameAlias       string
	keywordLenName         string
	keywordLenNameAlias    string
	keywordNumName         string
	keywordNumNameAlias    string
	keywordAmountName      string
	keywordAmountNameAlias string
}

func NewRuleItems(rule Ruler, opts ...func(items *RuleItems)) *RuleItems {
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

	ri.initName()
	ri.initLabels()
	ri.initTags()

	return ri
}

func (ri *RuleItems) initName() {
	ri.tableName = ri.rule.TableName()
	ri.name = ri.rule.GetName()
	ri.nameAlias = ri.genAlias(ri.name)

	ri.labelNumName = ri.rule.LabelNumName()
	ri.labelNumNameAlias = ri.genAlias(ri.labelNumName)

	ri.keywordName = ri.rule.KeywordName()
	ri.keywordNameAlias = ri.genAlias(ri.keywordName)
	ri.keywordLenName = ri.rule.KeywordLenName()
	ri.keywordLenNameAlias = ri.genAlias(ri.keywordLenName)
	ri.keywordNumName = ri.rule.KeywordNumName()
	ri.keywordNumNameAlias = ri.genAlias(ri.keywordNumName)
	ri.keywordAmountName = ri.rule.KeywordAmountName()
	ri.keywordAmountNameAlias = ri.genAlias(ri.keywordAmountName)
}

func (ri *RuleItems) initLabels() {
	var labels []string
	for label := range ri.rule.GetLabels() {
		labels = append(labels, label)
	}

	sort.Slice(labels, func(i, j int) bool {
		return labels[i] < labels[j]
	})

	ri.labels = labels

	var labelsAlias []string
	for _, label := range ri.labels {
		labelAlias, _ := ri.getAlias(label)

		labelsAlias = append(labelsAlias, labelAlias)
	}

	ri.labelsAlias = labelsAlias
}

func (ri *RuleItems) initTags() {
	ri.tags = []string{ri.name}
	ri.tags = append(ri.tags, ri.labels...)

	ri.tagsAlias = []string{ri.nameAlias}
	ri.tagsAlias = append(ri.tagsAlias, ri.labelsAlias...)
}

func (ri *RuleItems) getAlias(s string) (string, bool) {
	if ns, ok := ri.alias[s]; ok {
		return ns, ok
	} else {
		return s, ok
	}
}

func (ri *RuleItems) genAlias(s string) string {
	ns, _ := ri.getAlias(s)

	return ns
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
	return ri.name
}

func (ri *RuleItems) NameAlias() string {
	return ri.nameAlias
}

func (ri *RuleItems) TableName() string {
	return ri.tableName
}

func (ri *RuleItems) Labels() []string {
	return ri.labels
}

func (ri *RuleItems) LabelsAlias() []string {
	return ri.labelsAlias
}

func (ri *RuleItems) Tags() []string {
	return ri.tags
}

func (ri *RuleItems) TagsAlias() []string {
	return ri.tagsAlias
}

func (ri *RuleItems) LabelNumName() string {
	return ri.labelNumName
}

func (ri *RuleItems) LabelNumNameAlias() string {
	return ri.labelNumNameAlias
}

func (ri *RuleItems) KeywordName() string {
	return ri.keywordName
}

func (ri *RuleItems) KeywordNameAlias() string {
	return ri.keywordNameAlias
}

func (ri *RuleItems) KeywordLenName() string {
	return ri.keywordLenName
}

func (ri *RuleItems) KeywordLenNameAlias() string {
	return ri.keywordLenNameAlias
}

func (ri *RuleItems) KeywordNumName() string {
	return ri.keywordNumName
}

func (ri *RuleItems) KeywordNumNameAlias() string {
	return ri.keywordNumNameAlias
}

func (ri *RuleItems) KeywordAmountName() string {
	return ri.keywordAmountName
}

func (ri *RuleItems) KeywordAmountNameAlias() string {
	return ri.keywordAmountNameAlias
}
