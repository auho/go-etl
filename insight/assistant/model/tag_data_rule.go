package model

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ assistant.Moder = (*TagDataRule)(nil)

type TagDataRule struct {
	model
	data assistant.Dataor
	rule assistant.Ruler
}

func NewTagDataRule(data assistant.Dataor, rule assistant.Ruler, db *simpleDb.SimpleDB) *TagDataRule {
	t := &TagDataRule{}
	t.data = data
	t.rule = rule
	t.db = db

	return t
}

func (t *TagDataRule) GetData() assistant.Dataor {
	return t.data
}

func (t *TagDataRule) GetRule() assistant.Ruler {
	return t.rule
}

func (t *TagDataRule) GetName() string {
	return fmt.Sprintf("%s_%s", t.data.GetName(), t.rule.GetName())
}

func (t *TagDataRule) CommandExec(command *tablestructure.Command) {
	t.execCommand(command)
}

func (t *TagDataRule) GetDB() *simpleDb.SimpleDB {
	return t.db
}

func (t *TagDataRule) GetIdName() string {
	return "id"
}

func (t *TagDataRule) TableName() string {
	return fmt.Sprintf("%s_%s_%s", NameTag, t.data.GetName(), t.rule.GetName())
}
