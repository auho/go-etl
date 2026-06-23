package model

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpledb "github.com/auho/go-simple-db/v2"
)

var _ assistant.Entity = (*TagDataRule)(nil)

type TagDataRule struct {
	model
	extra
	data assistant.Entity
	rule assistant.Rule
}

func NewTagDataRule(data assistant.Entity, rule assistant.Rule, db *simpledb.SimpleDB) *TagDataRule {
	t := &TagDataRule{}
	t.data = data
	t.rule = rule
	t.db = db
	t.extra = extra{
		model: t,
	}

	return t
}

func (t *TagDataRule) GetData() assistant.Entity {
	return t.data
}

func (t *TagDataRule) GetRule() assistant.Rule {
	return t.rule
}

func (t *TagDataRule) GetName() string {
	return fmt.Sprintf("%s_%s", t.data.GetName(), t.rule.GetName())
}

func (t *TagDataRule) GetDB() *simpledb.SimpleDB {
	return t.db
}

func (t *TagDataRule) GetIDName() string {
	return "id"
}

func (t *TagDataRule) TableName() string {
	return fmt.Sprintf("%s_%s_%s", NameTag, t.data.GetName(), t.rule.GetName())
}

func (t *TagDataRule) WithCommand(fn func(*tablestructure.Command)) *TagDataRule {
	t.withCommand(fn)

	return t
}
