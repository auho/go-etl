package rule

import (
	"fmt"
	"strings"

	goetl "github.com/auho/go-etl"
	"github.com/auho/go-etl/tool"
	goSimpleDb "github.com/auho/go-simple-db/v2"
)

var defaultExcludeTableFields = []string{"id", "keyword_len"}

// WithDBRuleDataName
// rule prefix + data name + rule key name
func WithDBRuleDataName(n string) func(*DBRule) {
	return func(d *DBRule) {
		d.dataName = n
	}
}

// WithDBRuleShortTableName
// rule prefix + short table name
func WithDBRuleShortTableName(n string) func(*DBRule) {
	return func(d *DBRule) {
		d.shortTableName = n
	}
}

// WithDBRuleTagsName
// [tagA1, tagA2]
func WithDBRuleTagsName(s []string) func(*DBRule) {
	return func(r *DBRule) {
		r.tagsName = s
	}
}

func WithDBRuleExcludeFields(e []string) func(rule *DBRule) {
	return func(r *DBRule) {
		r.excludeFields = e
	}
}

// DBRule
// db rule
//
// table name
// - rule prefix + data name + rule key name
// - rule prefix + short table name
// - rule prefix + rule key name
type DBRule struct {
	db             *goSimpleDb.SimpleDB
	key            string   //
	tableName      string   // rule 表：rule_tagA
	dataName       string   // data name
	shortTableName string   // short name of tag data table
	tagsName       []string // 关键词匹配的标签列表名称： [tagA1, tagA2]
	excludeFields  []string
}

func NewDBRule(db *goSimpleDb.SimpleDB, key string, opts ...func(*DBRule)) *DBRule {
	d := &DBRule{
		db:  db,
		key: key,
	}

	for _, option := range opts {
		option(d)
	}

	if d.shortTableName != "" {
		d.tableName = goetl.RuleTableNamePrefix + "_" + d.shortTableName
	} else if d.dataName != "" {
		d.tableName = goetl.RuleTableNamePrefix + "_" + d.dataName + "_" + d.key
	} else {
		d.tableName = goetl.RuleTableNamePrefix + "_" + d.key
	}

	if len(d.excludeFields) <= 0 {
		d.excludeFields = defaultExcludeTableFields
	}

	return d
}

func (d *DBRule) TagsName() []string {
	if len(d.tagsName) <= 0 {
		row, err := d.db.GetTableColumns(d.tableName)
		if err != nil {
			panic(err)
		}

		_ef := make(map[string]struct{})
		for _, v := range d.excludeFields {
			_ef[v] = struct{}{}
		}

		for _, r := range row {
			if _, ok := _ef[r]; ok {
				continue
			}

			d.tagsName = append(d.tagsName, r)
		}
	} else {
		d.tagsName = append(d.tagsName, d.key+"_keyword")
	}

	d.tagsName = tool.RemoveReplicaSliceString(d.tagsName)

	return d.tagsName
}

func (d *DBRule) Items() []map[string]string {
	rules := make([]map[string]interface{}, 0)
	query := fmt.Sprintf("SELECT %s FROM %s ORDER BY keyword_len DESC, id ASC", strings.Join(d.tagsName, ", "), d.tableName)
	err := d.db.Raw(query).Scan(&rules).Error
	if err != nil {
		panic(err)
	}

	if len(rules) <= 0 {
		panic("rules is null")
	}

	_rules := make([]map[string]string, 0)
	for _, rule := range rules {
		_rule := make(map[string]string)
		for k, v := range rule {
			if _v, ok := v.(string); ok {
				_rule[k] = _v
			} else {
				panic(fmt.Sprintf("TagMatcher type of v is not string %s => %v", k, v))
			}
		}

		_rules = append(_rules, _rule)
	}

	return _rules
}
