package model

import (
	"fmt"

	simpleDb "github.com/auho/go-simple-db/v2"
)

type TagDataRule struct {
	data Dataor
	rule Ruler
	db   *simpleDb.SimpleDB
}

func NewTagDataRule(data Dataor, rule Ruler, db *simpleDb.SimpleDB) *TagDataRule {
	t := &TagDataRule{}
	t.data = data
	t.rule = rule
	t.db = db

	return t
}

func (t *TagDataRule) GetData() Dataor {
	return t.data
}

func (t *TagDataRule) GetRule() Ruler {
	return t.rule
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
