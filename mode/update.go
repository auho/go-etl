package mode

import (
	"github.com/auho/go-etl/means"
)

type UpdateMode struct {
	Mode
	update means.UpdateMeans
}

func NewTagUpdate(keys []string, update means.UpdateMeans) *UpdateMode {
	m := &UpdateMode{}
	m.keys = keys
	m.update = update

	return m
}

func (m *UpdateMode) GetFields() []string {
	return m.keys
}

func (m *UpdateMode) Do(item map[string]interface{}) map[string]interface{} {
	if item == nil {
		return nil
	}

	contents := m.GetKeysContent(m.keys, item)

	return m.update.Update(contents)
}

func (m *UpdateMode) Close() {
	m.update.Close()
}
