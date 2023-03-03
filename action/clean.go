package action

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/mode"
	go_simple_db "github.com/auho/go-simple-db/v2"
)

var _ Actionor = (*Clean)(nil)

type Clean struct {
	action
	target         *go_simple_db.SimpleDB
	fields         []string
	modes          []mode.UpdateModer
	targetDataName string
}

func NewClean(db *go_simple_db.SimpleDB, targetDataName string, modes []mode.UpdateModer) *Clean {
	var err error

	c := &Clean{}
	c.target = db
	c.modes = modes
	c.targetDataName = targetDataName

	c.fields, err = db.GetTableColumns(c.targetDataName)
	if err != nil {
		panic(err)
	}

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

	return fmt.Sprintf("Clean[%s] {%s}", c.targetDataName, strings.Join(s, ", "))
}

func (c *Clean) Prepare() error {
	return nil
}

func (c *Clean) Do(items []map[string]any) {
	targetItems := make([]map[string]any, 0)

	for _, item := range items {
		isClean := false
		for _, m := range c.modes {
			_res := m.Do(item)
			if _res != nil || len(_res) > 0 {
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
		targetItems = append(targetItems, item)
	}

	_ = c.target.BulkInsertFromSliceMap(c.targetDataName, targetItems, 2000)
}

func (c *Clean) AfterDo() {}
