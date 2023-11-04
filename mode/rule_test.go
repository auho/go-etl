package mode

import (
	"fmt"

	"github.com/auho/go-etl/v2/means/tag"
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
	return _ruleTableName
}

func (r *ruleTest) KeywordName() string {
	return "a"
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
	return []string{"b"}
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
	var rows []map[string]string
	err := db.Raw(fmt.Sprintf("SELECT `a`, `ab`, `a_keyword` FROM %s", _ruleTableName)).Scan(&rows).Error
	return rows, err
}
