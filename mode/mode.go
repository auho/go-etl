package mode

import (
	"fmt"
	"strconv"
	"strings"
)

type VoidModer interface {
	GetTitle() string
	Prepare() error
	Do(map[string]any)
	Close() error
}

type InsertModer interface {
	GetTitle() string
	GetKeys() []string   // 处理后的 key name
	GetFields() []string // source data 里的 key name
	Prepare() error
	Do(map[string]any) []map[string]any
	Close() error
}

type UpdateModer interface {
	GetTitle() string
	GetFields() []string // source data 里的 key name
	Prepare() error
	Do(map[string]any) map[string]any
	Close() error
}

type Mode struct {
	keys []string // 要被处理的 key name
}

func (t *Mode) getTitle() string {
	return "keys[" + strings.Join(t.keys, ", ") + "]"
}

func (t *Mode) GetKeysContent(keys []string, item map[string]any) []string {
	contents := make([]string, 0)
	for _, key := range keys {
		keyValue := t.KeyValueToString(key, item)

		contents = append(contents, keyValue)
	}

	return contents
}

func (t *Mode) KeyValueToString(key string, item map[string]any) string {
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
