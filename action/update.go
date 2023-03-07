package action

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/mode"
	"github.com/auho/go-etl/tool"
	goSimpleDb "github.com/auho/go-simple-db/v2"
)

var _ Actor = (*Update)(nil)

type Update struct {
	action
	db        *goSimpleDb.SimpleDB
	modes     []mode.UpdateModer
	idName    string
	dataTable string

	isTransfer    bool
	transferTable string
}

func NewUpdateAndTransfer(db *goSimpleDb.SimpleDB, dataTable string, transferTable string, idName string, modes []mode.UpdateModer) *Update {
	u := NewUpdate(db, dataTable, idName, modes)
	u.isTransfer = true
	u.transferTable = transferTable

	return u
}

func NewUpdate(db *goSimpleDb.SimpleDB, dataTable string, idName string, modes []mode.UpdateModer) *Update {
	u := &Update{}
	u.db = db
	u.modes = modes
	u.idName = idName
	u.dataTable = dataTable

	u.Init()

	return u
}

func (u *Update) GetFields() []string {
	fields := make([]string, 0)
	fields = append(fields, u.idName)

	for _, m := range u.modes {
		fields = append(fields, m.GetFields()...)
	}

	if u.isTransfer {
		columns, err := u.db.GetTableColumns(u.transferTable)
		if err != nil {
			panic(err)
		}

		fields = append(fields, columns...)
	}

	fields = tool.RemoveReplicaSliceString(fields)

	return fields
}

func (u *Update) Title() string {
	s := make([]string, 0)
	for _, m := range u.modes {
		s = append(s, m.GetTitle())
	}

	return fmt.Sprintf("Update[%s] {%s}", u.dataTable, strings.Join(s, ", "))
}

func (u *Update) Prepare() error {
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
		newItem[u.idName] = item[u.idName]
	}

	for k, v := range _does {
		newItem[k] = v
	}

	return []map[string]any{newItem}, true
}

func (u *Update) PostBatchDo(items []map[string]any) {
	var err error
	if u.isTransfer {
		err = u.db.BulkInsertFromSliceMap(u.transferTable, items, batchSize)
	} else {
		err = u.db.BulkUpdateFromSliceMapById(u.dataTable, u.idName, items)
	}

	if err != nil {
		panic(err)
	}
}

func (u *Update) PostDo() {}
