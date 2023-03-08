package action

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/mode"
	goSimpleDb "github.com/auho/go-simple-db/v2"
)

var _ Actor = (*Clean)(nil)

type Clean struct {
	action
	db          *goSimpleDb.SimpleDB
	fields      []string
	modes       []mode.UpdateModer
	targetTable string
}

func NewClean(db *goSimpleDb.SimpleDB, targetTable string, modes []mode.UpdateModer) *Clean {
	var err error

	c := &Clean{}
	c.db = db
	c.modes = modes
	c.targetTable = targetTable

	c.fields, err = db.GetTableColumns(c.targetTable)
	if err != nil {
		panic(err)
	}

	c.Init()

	return c
}

func (c *Clean) GetFields() []string {
	return c.fields
}

func (c *Clean) Title() string {
	s := make([]string, 0)
	for _, m := range c.modes {
		s = append(s, m.GetTitle())
	}

	return fmt.Sprintf("Clean[%s] {%s}", c.targetTable, strings.Join(s, ", "))
}

func (c *Clean) Prepare() error {
	return nil
}

func (c *Clean) Do(item map[string]any) ([]map[string]any, bool) {
	isClean := false
	for _, m := range c.modes {
		_res := m.Do(item)
		if len(_res) > 0 {
			isClean = true
			break
		}
	}

	if isClean == true {
		return nil, false
	}

	return []map[string]any{item}, true
}

func (c *Clean) PostBatchDo(items []map[string]any) {
	err := c.db.BulkInsertFromSliceMap(c.targetTable, items, batchSize)
	if err != nil {
		panic(err)
	}
}

func (c *Clean) PostDo() {}
