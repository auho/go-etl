// Package testutil provides shared test helpers for extract sub-packages.
package testutil

import "github.com/auho/go-etl/v3/task/extract"

// RuleTest implements extract.Rule with fixed test data.
type RuleTest struct{}

func (r *RuleTest) LabelNumName() string               { return "a_label_num" }
func (r *RuleTest) LabelNumNameAlias() string           { return "a_label_num" }
func (r *RuleTest) Name() string                        { return "a" }
func (r *RuleTest) NameAlias() string                   { return r.Name() }
func (r *RuleTest) TableName() string                   { return "rule_test" }
func (r *RuleTest) KeywordName() string                 { return "a_keyword" }
func (r *RuleTest) KeywordNameAlias() string             { return r.KeywordName() }
func (r *RuleTest) KeywordNumName() string               { return "a_keyword_num" }
func (r *RuleTest) KeywordNumNameAlias() string           { return r.KeywordNumName() }
func (r *RuleTest) KeywordAmountName() string             { return "a_keyword_amount" }
func (r *RuleTest) KeywordAmountNameAlias() string         { return r.KeywordAmountName() }
func (r *RuleTest) Labels() []string                     { return []string{"ab"} }
func (r *RuleTest) LabelsAlias() []string                 { return r.Labels() }
func (r *RuleTest) Tags() []string {
	return append([]string{r.Name()}, r.Labels()...)
}
func (r *RuleTest) TagsAlias() []string {
	return append([]string{r.NameAlias()}, r.LabelsAlias()...)
}
func (r *RuleTest) Fixed() map[string]string             { return nil }
func (r *RuleTest) FixedAlias() map[string]string         { return nil }
func (r *RuleTest) FixedKeys() []string                  { return nil }
func (r *RuleTest) FixedKeysAlias() []string              { return nil }
func (r *RuleTest) Items() ([]map[string]string, error) {
	return []map[string]string{
		{"a": "123", "ab": "123", "a_keyword": "123"},
		{"a": "a", "ab": "a1", "a_keyword": "b"},
		{"a": "e", "ab": "e1", "a_keyword": "e"},
		{"a": "中文", "ab": "中文1", "a_keyword": "中文"},
		{"a": "中1文", "ab": "中1文1", "a_keyword": `中_文`},
	}, nil
}
func (r *RuleTest) ItemsAlias() ([]map[string]string, error) {
	return r.Items()
}
func (r *RuleTest) ItemsForRegexp() ([]map[string]string, error) {
	return r.ItemsAlias()
}

// RuleAliasFixedTest embeds RuleTest and overrides methods for alias+fixed testing.
type RuleAliasFixedTest struct {
	RuleTest
}

func (r *RuleAliasFixedTest) NameAlias() string {
	return r.Name() + "_alias"
}
func (r *RuleAliasFixedTest) KeywordNameAlias() string {
	return r.KeywordName() + "_alias"
}
func (r *RuleAliasFixedTest) KeywordNumNameAlias() string {
	return r.KeywordNumName() + "_alias"
}
func (r *RuleAliasFixedTest) KeywordAmountNameAlias() string {
	return r.KeywordAmountName() + "_alias"
}
func (r *RuleAliasFixedTest) LabelsAlias() []string {
	var labels []string
	for _, label := range r.Labels() {
		labels = append(labels, label+"_alias")
	}
	return labels
}
func (r *RuleAliasFixedTest) Tags() []string {
	return append(append([]string{r.Name()}, r.Labels()...), r.FixedKeys()...)
}
func (r *RuleAliasFixedTest) TagsAlias() []string {
	return append(append([]string{r.NameAlias()}, r.LabelsAlias()...), r.FixedKeysAlias()...)
}
func (r *RuleAliasFixedTest) Fixed() map[string]string {
	return map[string]string{
		"c": "c_fixed",
		"d": "d_fixed",
	}
}
func (r *RuleAliasFixedTest) FixedAlias() map[string]string {
	return map[string]string{
		"c_alias": "c_fixed",
		"d_alias": "d_fixed",
	}
}
func (r *RuleAliasFixedTest) FixedKeys() []string {
	return []string{"c", "d"}
}
func (r *RuleAliasFixedTest) FixedKeysAlias() []string {
	return []string{"c_alias", "d_alias"}
}
func (r *RuleAliasFixedTest) ItemsAlias() ([]map[string]string, error) {
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
func (r *RuleAliasFixedTest) ItemsForRegexp() ([]map[string]string, error) {
	return r.ItemsAlias()
}

var _ extract.Rule = (*RuleTest)(nil)
var _ extract.Rule = (*RuleAliasFixedTest)(nil)