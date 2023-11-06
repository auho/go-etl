package model

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant"
	simpleDb "github.com/auho/go-simple-db/v2"
)

type TagDataRules struct {
	name  string
	data  assistant.Dataor
	rules []assistant.Ruler
	db    *simpleDb.SimpleDB
}

func NewTagDataSpreadRules(name string, data assistant.Dataor, rules []assistant.Ruler, db *simpleDb.SimpleDB) *TagDataRules {
	t := &TagDataRules{}
	t.name = name
	t.data = data
	t.rules = rules
	t.db = db

	return t
}

func (t *TagDataRules) GetData() assistant.Dataor {
	return t.data
}

func (t *TagDataRules) GetRules() []assistant.Ruler {
	return t.rules
}

func (t *TagDataRules) GetDB() *simpleDb.SimpleDB {
	return t.db
}

func (t *TagDataRules) GetIdName() string {
	return "id"
}

func (t *TagDataRules) TableName() string {
	return fmt.Sprintf("%s_%s_%s", NameTag, t.data.GetName(), t.name)
}
