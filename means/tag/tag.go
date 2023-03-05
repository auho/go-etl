package tag

import (
	"github.com/auho/go-etl/means/tag/rule"
	go_simple_db "github.com/auho/go-simple-db/v2"
)

// #TODO 重命名（key keyword  label=>tag） SQL use orm or interface

// Result
// result 匹配结果
type Result struct {
	Key        string            // keyword
	Num        int64             // matched num
	Texts      map[string]int64  // matched text map[matched text]num
	Tags       map[string]string // tags map[tag name]tag
	IsKeyMerge bool
}

func NewResult() *Result {
	m := &Result{}
	m.Tags = make(map[string]string)
	m.Texts = make(map[string]int64)

	return m
}

// LabelResult
// label result
type LabelResult struct {
	Identity string
	Labels   map[string]string         // tags map[tag name]tag
	Match    map[string]map[string]int // keyword and match text map[keyword]map[matched text]num
	KeyNum   int64                     // keyword num
	TextNum  int64                     // all text word num
}

func NewLabelResult() *LabelResult {
	l := &LabelResult{}
	l.Labels = make(map[string]string)
	l.Match = make(map[string]map[string]int)

	return l
}

type RuleMatcherOption func(*ruleMatcher)

// WithAlias
// alias [data name => output name]
func WithAlias(alias map[string]string) RuleMatcherOption {
	return func(r *ruleMatcher) {
		r.alias = alias
	}
}

// WithFixedTags
// fixed data map[tagName]tag
func WithFixedTags(m map[string]string) RuleMatcherOption {
	return func(r *ruleMatcher) {
		r.fixedTags = m
	}
}

// WithRule
// rule
func WithRule(rule rule.Ruler) RuleMatcherOption {
	return func(r *ruleMatcher) {
		r.rule = rule
	}
}

// WithDBRule
// with db rule
func WithDBRule(db *go_simple_db.SimpleDB, opts ...func(*rule.DBRule)) RuleMatcherOption {
	return func(r *ruleMatcher) {
		r.postOptionsFuncs = append(r.postOptionsFuncs, func(rr *ruleMatcher) {
			fn := WithRule(rule.NewDBRule(db, rr.key, opts...))
			fn(rr)
		})
	}
}

// WithMatcher
// matcher
func WithMatcher(m *Matcher) RuleMatcherOption {
	return func(r *ruleMatcher) {
		r.matcher = m
	}
}

// WithMatcherOptions
// with matcher options
func WithMatcherOptions(opts ...MatcherOption) RuleMatcherOption {
	return func(r *ruleMatcher) {
		r.postOptionsFuncs = append(r.postOptionsFuncs, WithMatcher(NewMatcher(opts...)))
	}
}

// ruleMatcher
// tag matcher
type ruleMatcher struct {
	matcher         *Matcher
	rule            rule.Ruler        // rule
	key             string            // 关键词名称： tagA
	keyFieldName    string            // 关键词字段名称：tagA_keyword
	keyNumFieldName string            // 关键词出现数量：tagA_keyword_num
	tagsName        []string          // 关键词匹配的标签列表名称： [tagA1, tagA2]
	alias           map[string]string // alias [data name => output name]
	fixedTags       map[string]string // fixed tags data
	fixedKeys       []string          // keys of fixed data

	postOptionsFuncs []RuleMatcherOption
}

func newRuleMatcher(key string, opts ...RuleMatcherOption) *ruleMatcher {
	r := &ruleMatcher{}
	r.key = key

	for _, option := range opts {
		option(r)
	}

	for _, p := range r.postOptionsFuncs {
		p(r)
	}

	if r.rule == nil {
		panic("rule is nil")
	}

	r.prepare()

	return r
}

func (r *ruleMatcher) prepare() {
	r.keyFieldName = r.key + "_keyword"
	r.keyNumFieldName = r.keyFieldName + "_num"

	// 获取被匹配的标签字段列表
	tags := r.rule.TagsName()
	for _, v := range tags {
		if v != r.keyFieldName {
			r.tagsName = append(r.tagsName, v)
		}
	}

	rules := r.rule.Items()
	if len(r.alias) > 0 {
		for k := range rules {
			for ak, av := range r.alias {
				if rv, ok := rules[k][ak]; ok {
					rules[k][av] = rv
					delete(rules[k], ak)
				}
			}
		}
	}

	r.prepareAlias()

	if r.matcher == nil {
		r.matcher = DefaultMatcher
	}

	r.matcher.prepare(r.keyFieldName, rules)
}

func (r *ruleMatcher) prepareAlias() {
	if len(r.alias) <= 0 {
		return
	}

	for k, v := range r.tagsName {
		if s, ok := r.alias[v]; ok {
			r.tagsName[k] = s
		}
	}

	if s, ok := r.alias[r.keyFieldName]; ok {
		r.keyFieldName = s
	}

	if s, ok := r.alias[r.keyNumFieldName]; ok {
		r.keyNumFieldName = s
	}

	r.fixedKeys = make([]string, 0)
	for k, v := range r.fixedTags {
		if av, ok := r.alias[k]; ok {
			r.fixedTags[av] = v
			delete(r.fixedTags, k)
		}
	}
}

func (r *ruleMatcher) getResultInsertKeys() []string {
	return append([]string{r.keyFieldName, r.keyNumFieldName}, append(r.tagsName, r.fixedKeys...)...)
}

func (r *ruleMatcher) resultsToSliceMap(results []*Result) []map[string]interface{} {
	items := make([]map[string]interface{}, 0, len(results))
	for _, result := range results {
		items = append(items, r.resultToMap(result))
	}

	return items
}

func (r *ruleMatcher) resultToSliceMap(result *Result) []map[string]interface{} {
	return []map[string]interface{}{r.resultToMap(result)}
}

func (r *ruleMatcher) resultToMap(result *Result) map[string]interface{} {
	item := make(map[string]interface{})
	item[r.keyFieldName] = result.Key
	item[r.keyNumFieldName] = result.Num

	for _, tagName := range r.tagsName {
		item[tagName] = result.Tags[tagName]
	}

	if len(r.fixedTags) > 0 {
		for k := range r.fixedTags {
			item[k] = r.fixedTags[k]
		}
	}

	return item
}

// Results
// results
type Results []*Result

func (r Results) Len() int {
	return len(r)
}

func (r Results) Less(i, j int) bool {
	return r[i].Num > r[j].Num
}

func (r Results) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// LabelResults
// label results
type LabelResults []*LabelResult

func (l LabelResults) Len() int {
	return len(l)
}

func (l LabelResults) Less(i, j int) bool {
	return l[i].TextNum > l[j].TextNum
}

func (l LabelResults) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
