package tag

import (
	"github.com/auho/go-etl/v2/job/means"
)

var _ means.Ruler = (*ruleTest)(nil)
var _ means.Ruler = (*ruleAliasFixedTest)(nil)

// ruleTest
// rule
type ruleTest struct {
}

func (r *ruleTest) Name() string {
	return "a"
}

func (r *ruleTest) NameAlias() string {
	return r.Name()
}

func (r *ruleTest) TableName() string {
	return _ruleTableName
}

func (r *ruleTest) KeywordName() string {
	return "a_keyword"
}

func (r *ruleTest) KeywordNameAlias() string {
	return r.KeywordName()
}

func (r *ruleTest) KeywordNumName() string {
	return "a_keyword_num"
}

func (r *ruleTest) KeywordNumNameAlias() string {
	return r.KeywordNumName()
}

func (r *ruleTest) Labels() []string {
	return []string{"ab"}
}

func (r *ruleTest) LabelsAlias() []string {
	return r.Labels()
}

func (r *ruleTest) Fixed() map[string]any {
	return nil
}

func (r *ruleTest) FixedAlias() map[string]any {
	return nil
}

func (r *ruleTest) FixedKeys() []string {
	return nil
}

func (r *ruleTest) FixedKeysAlias() []string {
	return nil
}

func (r *ruleTest) Items() ([]map[string]string, error) {
	return []map[string]string{
		{"a": "a", "ab": "a1", "a_keyword": "a"},
		{"a": "a", "ab": "a1", "a_keyword": "b"},
		{"a": "ab", "ab": "ab1", "a_keyword": "ab"},
		{"a": "123", "ab": "123", "a_keyword": "123"},
		{"a": "中文", "ab": "中文1", "a_keyword": "中文"},
	}, nil
}

func (r *ruleTest) ItemsAlias() ([]map[string]string, error) {
	return r.Items()
}

func (r *ruleTest) ItemsForRegexp() ([]map[string]string, error) {
	return r.ItemsAlias()
}

// ruleAliasFixedTest
// rule alias fixed
type ruleAliasFixedTest struct {
	ruleTest
}

func (r *ruleAliasFixedTest) NameAlias() string {
	return r.Name() + "_alias"
}

func (r *ruleAliasFixedTest) KeywordNameAlias() string {
	return r.KeywordName() + "_alias"
}

func (r *ruleAliasFixedTest) KeywordNumNameAlias() string {
	return r.KeywordNumName() + "_alias"
}

func (r *ruleAliasFixedTest) LabelsAlias() []string {
	var labels []string
	for _, label := range r.Labels() {
		labels = append(labels, label+"_alias")
	}

	return labels
}

func (r *ruleAliasFixedTest) Fixed() map[string]any {
	return map[string]any{
		"c": "c_fixed",
		"d": "d_fixed",
	}
}

func (r *ruleAliasFixedTest) FixedAlias() map[string]any {
	return map[string]any{
		"c_alias": "c_fixed",
		"d_alias": "d_fixed",
	}
}

func (r *ruleAliasFixedTest) FixedKeys() []string {
	return []string{"c", "d"}
}

func (r *ruleAliasFixedTest) FixedKeysAlias() []string {
	return []string{"c_alias", "d_alias"}
}

func (r *ruleAliasFixedTest) ItemsAlias() ([]map[string]string, error) {
	var newItems []map[string]string

	items, _ := r.Items()
	for _, v := range items {
		nv := make(map[string]string, len(v))
		for _k, _v := range v {
			nv[_k+"_alias"] = _v
		}

		newItems = append(newItems, nv)
	}

	return newItems, nil
}

func (r *ruleAliasFixedTest) ItemsForRegexp() ([]map[string]string, error) {
	return r.ItemsAlias()
}
