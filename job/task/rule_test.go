package task

import (
	"fmt"

	"github.com/auho/go-etl/v2/job/means/tag"
)

var _ tag.Ruler = (*ruleTest)(nil)

type ruleTest struct {
}

func (r *ruleTest) Name() string {
	return "a"
}

func (r *ruleTest) NameAlias() string {
	return r.Name()
}

func (r *ruleTest) TableName() string {
	return _ruleTable
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
	err := _db.Raw(fmt.Sprintf("SELECT `a`, `ab`, `a_keyword` FROM %s", _ruleTable)).Scan(&rows).Error

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
