package entity

import (
	"fmt"

	"github.com/auho/go-etl/v3/insight/assistant"
	"github.com/auho/go-etl/v3/insight/assistant/schema"
	simpledb "github.com/auho/go-simple-db/v3"
)

var _ assistant.Entity = (*TagDataRules)(nil)

type TagDataRules struct {
	model
	extra
	name  string
	data  assistant.Entity
	rules []assistant.Rule
}

func NewTagDataRules(name string, data assistant.Entity, rules []assistant.Rule, db *simpledb.SimpleDB) *TagDataRules {
	t := &TagDataRules{}
	t.name = name
	t.data = data
	t.rules = rules
	t.db = db
	t.extra = extra{
		model: t,
	}

	return t
}

func (t *TagDataRules) GetData() assistant.Entity {
	return t.data
}

func (t *TagDataRules) GetRules() []assistant.Rule {
	return t.rules
}

func (t *TagDataRules) GetName() string {
	return fmt.Sprintf("%s_%s", t.data.GetName(), t.name)
}

func (t *TagDataRules) DB() *simpledb.SimpleDB {
	return t.db
}

func (t *TagDataRules) IDName() string {
	return "id"
}

func (t *TagDataRules) TableName() string {
	return fmt.Sprintf("%s_%s_%s", NameTag, t.data.GetName(), t.name)
}

func (t *TagDataRules) WithCommand(fn func(*schema.Command)) *TagDataRules {
	t.withCommand(fn)

	return t
}

func (t *TagDataRules) Clone(name string) *TagDataRules {
	return NewTagDataRules(name, t.data, t.rules, t.db)
}
