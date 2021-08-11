package mode

import (
	"github.com/auho/go-etl/means"
)

type Update struct {
	Mode
	update means.UpdateMeans
}

func NewUpdate(keys []string, update means.UpdateMeans) *Update {
	m := &Update{}
	m.keys = keys
	m.update = update

	return m
}

func (m *Update) GetTitle() string {
	return m.getModeTitle() + " " + m.update.GetTitle()
}

func (m *Update) GetFields() []string {
	return m.keys
}

func (m *Update) Do(item map[string]interface{}) map[string]interface{} {
	if item == nil {
		return nil
	}

	contents := m.GetKeysContent(m.keys, item)

	return m.update.Update(contents)
}

func (m *Update) Close() {
	m.update.Close()
}
