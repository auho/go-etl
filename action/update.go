package action

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/mode"
	"github.com/auho/go-etl/tool"
	go_simple_db "github.com/auho/go-simple-db/v2"
)

var _ Actionor = (*Update)(nil)

type Update struct {
	action
	target         *go_simple_db.SimpleDB
	modes          []mode.UpdateModer
	targetDataName string
	idName         string
	dataName       string
}

func NewUpdate(db *go_simple_db.SimpleDB, dataName string, idName string, modes []mode.UpdateModer) *Update {
	u := &Update{}
	u.target = db
	u.modes = modes
	u.idName = idName
	u.dataName = dataName

	return u
}

func (u *Update) GetFields() []string {
	fields := make([]string, 0)
	for _, m := range u.modes {
		fields = append(fields, m.GetFields()...)
	}

	fields = tool.RemoveReplicaSliceString(fields)

	return append(fields, u.idName)
}

func (u *Update) Title() string {
	s := make([]string, 0)
	for _, m := range u.modes {
		s = append(s, m.GetTitle())
	}

	return fmt.Sprintf("Update[%s] {%s}", u.dataName, strings.Join(s, ", "))
}

func (u *Update) Prepare() error {
	//TODO implement me
	panic("implement me")
}

func (u *Update) Do(items []map[string]any) {
	targetItems := make([]map[string]interface{}, 0)

	for _, item := range items {
		_newItem := make(map[string]any)
		for _, m := range u.modes {
			_do := m.Do(item)
			for k, v := range _do {
				_newItem[k] = v
			}
		}

		if len(_newItem) <= 0 {
			continue
		}

		_newItem[u.idName] = item[u.idName]
		targetItems = append(targetItems, _newItem)
	}

	_ = u.target.BulkUpdateFromSliceMapById(u.dataName, u.idName, targetItems)

}

func (u *Update) AfterDo() {}
