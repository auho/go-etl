package mode

import (
	"fmt"

	"github.com/auho/go-etl/means"
	"github.com/auho/go-simple-db/simple"
)

type CleanMode struct {
	Mode
	fields   []string
	updates  []means.UpdateMeans
	dataName string
}

func NewCleanMode(dataName string, updates []means.UpdateMeans, db simple.Driver) *CleanMode {
	m := &CleanMode{}
	m.updates = updates
	m.dataName = dataName

	var err error
	m.fields, err = db.GetTableColumns(dataName)
	if err != nil {
		panic(err)
	}

	return m
}

func (m *CleanMode) GetTitle() string {
	return fmt.Sprintf("Clean[%s]", m.dataName)
}

func (m *CleanMode) GetFields() []string {
	return m.fields
}

func (m *CleanMode) GetKeys() []string {
	return m.fields
}

func (m *CleanMode) Do(item map[string]interface{}) [][]interface{} {
	if item == nil {
		return nil
	}

	contents := m.GetKeysContent(m.keys, item)

	result := make([]interface{}, 0)
	isClean := false
	for _, u := range m.updates {
		res := u.Update(contents)
		if res == nil {
			continue
		} else {
			isClean = true
			break
		}
	}

	if isClean {
		return nil
	} else {
		return [][]interface{}{result}
	}
}

func (m *CleanMode) Close() {
	for _, u := range m.updates {
		u.Close()
	}
}
