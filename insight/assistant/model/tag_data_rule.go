package model

import (
	"fmt"

	"github.com/auho/go-etl/v2/insight/assistant"
	"github.com/auho/go-etl/v2/insight/assistant/accessory/dml"
	"github.com/auho/go-etl/v2/insight/assistant/tablestructure"
	simpleDb "github.com/auho/go-simple-db/v2"
)

var _ assistant.Moder = (*TagDataRule)(nil)

type TagDataRule struct {
	model
	data assistant.Rowsor
	rule assistant.Ruler
}

func NewTagDataRule(data assistant.Rowsor, rule assistant.Ruler, db *simpleDb.SimpleDB) *TagDataRule {
	t := &TagDataRule{}
	t.data = data
	t.rule = rule
	t.db = db

	return t
}

func (t *TagDataRule) GetData() assistant.Rowsor {
	return t.data
}

func (t *TagDataRule) GetRule() assistant.Ruler {
	return t.rule
}

func (t *TagDataRule) GetName() string {
	return fmt.Sprintf("%s_%s", t.data.GetName(), t.rule.GetName())
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

func (t *TagDataRule) DmlTable() *dml.Table {
	return dml.NewTable(t.TableName())
}

func (t *TagDataRule) WithCommand(fn func(*tablestructure.Command)) *TagDataRule {
	t.withCommand(fn)

	return t
}

func (t *TagDataRule) CopyBuild(dst assistant.Rawer) error {
	return t.db.DropAndCopy(t.TableName(), dst.TableName())
}
