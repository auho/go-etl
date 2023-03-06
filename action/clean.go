package action

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/mode"
	goSimpleDb "github.com/auho/go-simple-db/v2"
)

var _ Actioner = (*Clean)(nil)

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

func (c *Clean) Do(items []map[string]any) {
	newItems := make([]map[string]any, 0)

	for _, item := range items {
		isClean := false
		for _, m := range c.modes {
			_res := m.Do(item)
			if len(_res) > 0 {
				isClean = true
				break
			}
		}

		if isClean == true {
			continue
		}

		nm := make(map[string]any, len(c.fields))
		for _, field := range c.fields {
			nm[field] = item[field]
		}

		c.AddAmount(1)
		newItems = append(newItems, item)
	}

	err := c.db.BulkInsertFromSliceMap(c.targetTable, newItems, batchSize)
	if err != nil {
		panic(err)
	}
}

func (c *Clean) AfterDo() {}
