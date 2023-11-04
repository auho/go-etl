package action

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/mode"
)

var _ Actor = (*Clean)(nil)

// Clean
// filter
type Clean struct {
	action
	target Target
	modes  []mode.UpdateModer
	keys   []string
}

func NewClean(target Target, modes []mode.UpdateModer) *Clean {
	c := &Clean{}
	c.target = target
	c.modes = modes

	c.Init()

	return c
}

func (c *Clean) GetFields() []string {
	var err error
	c.keys, err = c.target.GetDB().GetTableColumns(c.target.TableName())
	if err != nil {
		panic(fmt.Errorf("GetTableColumns error; %w", err))
	}

	return c.keys
}

func (c *Clean) Title() string {
	s := make([]string, 0)
	for _, m := range c.modes {
		s = append(s, m.GetTitle())
	}

	return fmt.Sprintf("Clean[%s] {%s}", c.target.TableName(), strings.Join(s, ", "))
}

func (c *Clean) Prepare() error {
	for _, m := range c.modes {
		err := m.Prepare()
		if err != nil {
			return fmt.Errorf("clean action prepare error; %w", err)
		}
	}

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
	err := c.target.GetDB().BulkInsertFromSliceMap(c.target.TableName(), items, batchSize)
	if err != nil {
		panic(err)
	}
}

func (c *Clean) PostDo() {}
