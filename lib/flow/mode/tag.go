package mode

import (
	"etl/lib/flow/means"
)

type TagInsert struct {
	keys   []string
	insert means.InsertMeans
}

func NewInsert(keys []string, insert means.InsertMeans) *TagInsert {
	t := &TagInsert{}
	t.keys = keys
	t.insert = insert

	return t
}

func (ti *TagInsert) GetFields() []string {
	return ti.keys
}

func (ti *TagInsert) GetKeys() []string {
	return ti.insert.GetKeys()
}

func (ti *TagInsert) Do(item map[string]interface{}) [][]interface{} {
	if item == nil {
		return nil
	}

	contents := make([]string, 0)
	for _, key := range ti.keys {
		contents = append(contents, key)
	}

	return ti.insert.Insert(contents)
}

type TagUpdate struct {
	keys   []string
	update means.UpdateMeans
}

func NewUpdate(keys []string, update means.UpdateMeans) *TagUpdate {
	t := &TagUpdate{}
	t.keys = keys
	t.update = update

	return t
}

func (tu *TagUpdate) GetFields() []string {
	return tu.keys
}

func (tu *TagUpdate) Do(item map[string]interface{}) map[string]interface{} {
	if item == nil {
		return nil
	}

	contents := make([]string, 0)
	for _, key := range tu.keys {
		contents = append(contents, key)
	}

	return tu.update.Update(contents)
}
