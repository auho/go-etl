package explore

import (
	"fmt"

	"github.com/auho/go-etl/v2/job/means"
)

var _ means.Ruler = (*ruleTest)(nil)

type ruleTest struct {
}

func (r *ruleTest) LabelNumName() string {
	//TODO implement
	panic("implement me")
}

func (r *ruleTest) LabelNumNameAlias() string {
	//TODO implement me
	panic("implement me")
}

func (r *ruleTest) MeansKeys() []string {
	var keys []string
	keys = []string{
		r.NameAlias(),
		r.KeywordNameAlias(),
		r.KeywordNumNameAlias(),
	}
	keys = append(keys, r.LabelsAlias()...)
	keys = append(keys, r.FixedKeysAlias()...)

	return keys
}

func (r *ruleTest) MeansDefaultValues() map[string]any {
	defaultValues := map[string]any{
		r.NameAlias():           "",
		r.KeywordNameAlias():    "",
		r.KeywordNumNameAlias(): 0,
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
	return "a_keyword_Amount"
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

func (r *ruleTest) ItemsAlias() ([]map[string]string, error) {
	var rows []map[string]any
	err := _db.Raw(fmt.Sprintf("SELECT `a`, `ab`, `a_keyword` FROM %s", _ruleTableName)).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("scan error; %w", err)
	}

	var _rows []map[string]string
	for _, row := range rows {
		_row := make(map[string]string)
		for k, v := range row {
			_v := v.(string)
			_row[k] = _v
		}

		_rows = append(_rows, _row)
	}

	return _rows, err
}

func (r *ruleTest) ItemsForRegexp() ([]map[string]string, error) {
	return r.ItemsAlias()
}
