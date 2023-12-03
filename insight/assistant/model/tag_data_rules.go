package model

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ assistant.Moder = (*TagDataRules)(nil)

type TagDataRules struct {
	model
	extra
	name  string
	data  assistant.Rowsor
	rules []assistant.Ruler
}

func NewTagDataRules(name string, data assistant.Rowsor, rules []assistant.Ruler, db *simpleDb.SimpleDB) *TagDataRules {
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

func (t *TagDataRules) GetData() assistant.Rowsor {
	return t.data
}

func (t *TagDataRules) GetRules() []assistant.Ruler {
	return t.rules
}

func (t *TagDataRules) GetName() string {
	return fmt.Sprintf("%s_%s", t.data.GetName(), t.name)
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

func (t *TagDataRules) DmlTable() *dml.Table {
	return dml.NewTable(t.TableName())
}

func (t *TagDataRules) WithCommand(fn func(*tablestructure.Command)) *TagDataRules {
	t.withCommand(fn)

	return t
}

func (t *TagDataRules) Clone(name string) *TagDataRules {
	return NewTagDataRules(name, t.data, t.rules, t.db)
}
