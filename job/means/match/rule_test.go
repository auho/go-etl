package match

import (
	"github.com/auho/go-etl/v2/job/means"
)

var _ means.Ruler = (*ruleTest)(nil)
var _ means.Ruler = (*ruleAliasFixedTest)(nil)

// ruleTest
// rule
type ruleTest struct {
}

func (r *ruleTest) LabelNumName() string {
	return "a_label_num"
}

func (r *ruleTest) LabelNumNameAlias() string {
	return "a_label_num"
}

func (r *ruleTest) MeansKeys() []string {
	var keys []string
	keys = []string{
		r.NameAlias(),
		r.KeywordNameAlias(),
		r.KeywordAmountNameAlias(),
	}
	keys = append(keys, r.LabelsAlias()...)
	keys = append(keys, r.FixedKeysAlias()...)

	return keys
}

func (r *ruleTest) MeansDefaultValues() map[string]any {
	defaultValues := map[string]any{
		r.NameAlias():              "",
		r.KeywordNameAlias():       "",
		r.KeywordAmountNameAlias(): 0,
	}

	for _, _la := range r.LabelsAlias() {
		defaultValues[_la] = ""
	}

	for _, _fka := range r.FixedKeysAlias() {
		defaultValues[_fka] = ""
	}

	return defaultValues
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

func (r *ruleTest) KeywordAmountName() string {
	return "a_keyword_amount"
}

func (r *ruleTest) KeywordAmountNameAlias() string {
	return r.KeywordAmountName()
}

func (r *ruleTest) Labels() []string {
	return []string{"ab"}
}

func (r *ruleTest) LabelsAlias() []string {
	return r.Labels()
}

func (r *ruleTest) Tags() []string {
	return append([]string{r.Name()}, r.Labels()...)
}

func (r *ruleTest) TagsAlias() []string {
	return append([]string{r.NameAlias()}, r.LabelsAlias()...)
}

func (r *ruleTest) Fixed() map[string]string {
	return nil
}

func (r *ruleTest) FixedAlias() map[string]string {
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
		{"a": "e", "ab": "e1", "a_keyword": "e"},
		{"a": "ab", "ab": "ab1", "a_keyword": "ab"},
		{"a": "123", "ab": "123", "a_keyword": "123"},
		{"a": "中文", "ab": "中文1", "a_keyword": "中文"},
		{"a": "中1文", "ab": "中1文1", "a_keyword": `中_文`},
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

func (r *ruleAliasFixedTest) KeywordAmountNameAlias() string {
	return r.KeywordAmountName() + "_alias"
}

func (r *ruleAliasFixedTest) LabelsAlias() []string {
	var labels []string
	for _, label := range r.Labels() {
		labels = append(labels, label+"_alias")
	}

	return labels
}
func (r *ruleAliasFixedTest) Tags() []string {
	return append(append([]string{r.Name()}, r.Labels()...), r.FixedKeys()...)
}

func (r *ruleAliasFixedTest) TagsAlias() []string {
	return append(append([]string{r.NameAlias()}, r.LabelsAlias()...), r.FixedKeysAlias()...)
}

func (r *ruleAliasFixedTest) Fixed() map[string]string {
	return map[string]string{
		"c": "c_fixed",
		"d": "d_fixed",
	}
}

func (r *ruleAliasFixedTest) FixedAlias() map[string]string {
	return map[string]string{
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

func (r *ruleAliasFixedTest) MeansKeys() []string {
	var keys []string
	keys = []string{
		r.NameAlias(),
		r.KeywordNameAlias(),
		r.KeywordAmountNameAlias(),
	}
	keys = append(keys, r.LabelsAlias()...)
	keys = append(keys, r.FixedKeysAlias()...)

	return keys
}

func (r *ruleAliasFixedTest) MeansDefaultValues() map[string]any {
	defaultValues := map[string]any{
		r.NameAlias():              "",
		r.KeywordNameAlias():       "",
		r.KeywordAmountNameAlias(): 0,
	}

	for _, _la := range r.LabelsAlias() {
		defaultValues[_la] = ""
	}

	for _, _fka := range r.FixedKeysAlias() {
		defaultValues[_fka] = ""
	}

	return defaultValues
}
