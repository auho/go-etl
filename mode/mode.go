package mode

import (
	"fmt"
	"strconv"
	"strings"
)

type VoidModer interface {
	GetTitle() string
	Do(map[string]interface{})
	Close()
}

type InsertModer interface {
	GetTitle() string
	GetKeys() []string
	GetFields() []string
	Do(map[string]interface{}) []map[string]interface{}
	Close()
}

type UpdateModer interface {
	GetTitle() string
	GetFields() []string
	Do(map[string]interface{}) map[string]interface{}
	Close()
}

type Mode struct {
	keys []string
}

func (t *Mode) getModeTitle() string {
	return "keys[" + strings.Join(t.keys, ", ") + "]"
}

func (t *Mode) GetKeysContent(keys []string, item map[string]interface{}) []string {
	contents := make([]string, 0)
	for _, key := range keys {
		keyValue := t.KeyValueToString(key, item)

		contents = append(contents, keyValue)
	}

	return contents
}

func (t *Mode) KeyValueToString(key string, item map[string]interface{}) string {
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

	return keyValue
}
