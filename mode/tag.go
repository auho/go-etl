package mode

import (
	"fmt"
	"strconv"

	"github.com/auho/go-etl/means"
)

type TagInsert struct {
	keys   []string
	insert means.InsertMeans
}

func NewTagInsert(keys []string, insert means.InsertMeans) *TagInsert {
	t := &TagInsert{}
	t.keys = keys
	t.insert = insert

	return t
}

func (ti *TagInsert) GetName() string {
	return ti.insert.GetName()
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
		keyValue := ""

		switch item[key].(type) {
		case string:
			keyValue = item[key].(string)
		case []uint8:
			keyValue = string(item[key].([]uint8))
		case int64:
			keyValue = strconv.FormatInt(item[key].(int64), 10)
		default:
			panic(fmt.Sprintf("type is not string %T", item[key]))
		}

		contents = append(contents, keyValue)
	}

	return ti.insert.Insert(contents)
}

type TagUpdate struct {
	keys   []string
	update means.UpdateMeans
}

func NewTagUpdate(keys []string, update means.UpdateMeans) *TagUpdate {
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
