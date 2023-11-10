package action

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/job"
	"github.com/auho/go-etl/v2/job/mode"
	"github.com/auho/go-etl/v2/tool/slices"
)

var _ Actor = (*Update)(nil)

type Update struct {
	targetAction

	source job.Source
	modes  []mode.UpdateModer

	isTransfer bool
}

func NewUpdateAndTransfer(source job.Source, target job.Target, modes []mode.UpdateModer) *Update {
	u := NewUpdate(source, modes)
	u.target = target
	u.isTransfer = true

	return u
}

func NewUpdate(source job.Source, modes []mode.UpdateModer) *Update {
	u := &Update{}
	u.source = source
	u.modes = modes

	u.Init()

	return u
}

func (u *Update) GetFields() []string {
	fields := make([]string, 0)
	fields = append(fields, u.source.GetIdName())

	for _, m := range u.modes {
		fields = append(fields, m.GetFields()...)
	}

	if u.isTransfer {
		columns, err := u.target.GetDB().GetTableColumns(u.target.TableName())
		if err != nil {
			panic(err)
		}

		fields = append(fields, columns...)
	}

	fields = slices.SliceDropDuplicates(fields)

	return fields
}

func (u *Update) Title() string {
	s := make([]string, 0)
	for _, m := range u.modes {
		s = append(s, m.GetTitle())
	}

	return fmt.Sprintf("Update[%s] {%s}", u.source.TableName(), strings.Join(s, ", "))
}

func (u *Update) Prepare() error {
	for _, m := range u.modes {
		err := m.Prepare()
		if err != nil {
			return fmt.Errorf("update action prepare error; %w", err)
		}
	}

	return nil
}

func (u *Update) Do(item map[string]any) ([]map[string]any, bool) {
	_does := make(map[string]any)
	for _, m := range u.modes {
		_do := m.Do(item)
		for k, v := range _do {
			_does[k] = v
		}
	}

	if len(_does) <= 0 && u.isTransfer == false {
		return nil, false
	}

	var newItem map[string]any
	if u.isTransfer {
		newItem = item
	} else {
		newItem = make(map[string]any)
		newItem[u.source.GetIdName()] = item[u.source.GetIdName()]
	}

	for k, v := range _does {
		newItem[k] = v
	}

	return []map[string]any{newItem}, true
}

func (u *Update) PostBatchDo(items []map[string]any) {
	var err error
	if u.isTransfer {
		err = u.target.GetDB().BulkInsertFromSliceMap(u.target.TableName(), items, batchSize)
	} else {
		err = u.source.GetDB().BulkUpdateFromSliceMapById(u.source.TableName(), u.source.GetIdName(), items)
	}

	if err != nil {
		panic(err)
	}
}

func (u *Update) Blink()        {}
func (u *Update) PreDo() error  { return nil }
func (u *Update) PostDo() error { return nil }
func (u *Update) Close() error  { return nil }
