package mode

import (
	"fmt"
	"strconv"
	"strings"
	"sync/atomic"
)

type Moder interface {
	GetTitle() string
	GetFields() []string // source data 里的 key name
	Prepare() error
	Close() error
}

type VoidModer interface {
	Moder
	Do(map[string]any) map[string]any
}

type InsertModer interface {
	Moder
	GetKeys() []string             // 处理后的 key name
	DefaultValues() map[string]any // 需要 implement clone important!
	Do(map[string]any) []map[string]any
	State() []string
}

type UpdateModer interface {
	Moder
	Do(map[string]any) map[string]any
}

type TransferModer interface {
	Moder
	Do(map[string]any) map[string]any
}

type Mode struct {
	Keys   []string // 要被处理的 key name
	total  int64
	amount int64
}

func (m *Mode) AddTotal(num int64) {
	atomic.AddInt64(&m.total, num)
}

func (m *Mode) AddAmount(num int64) {
	atomic.AddInt64(&m.amount, num)
}

func (m *Mode) GenCounter() string {
	return fmt.Sprintf("total: %d; amount: %d", m.total, m.amount)
}

func (m *Mode) GenTitle(name string, means string) string {
	return fmt.Sprintf("%s %s{%s}", name, "keys["+strings.Join(m.Keys, ", ")+"]", means)
}

func (m *Mode) GetKeyContent(key string, item map[string]any) string {
	return m.KeyValueToString(key, item)
}

func (m *Mode) GetKeysContent(keys []string, item map[string]any) []string {
	contents := make([]string, 0)
	for _, key := range keys {
		keyValue := m.KeyValueToString(key, item)

		contents = append(contents, keyValue)
	}

	return contents
}

func (m *Mode) KeyValueToString(key string, item map[string]any) string {
	keyValue := ""

	switch item[key].(type) {
	case string:
		keyValue = item[key].(string)
	case []uint8:
		keyValue = string(item[key].([]uint8))
	case int64:
		keyValue = strconv.FormatInt(item[key].(int64), 10)
	case nil:

	default:
		panic(fmt.Sprintf("type is not string %T", item[key]))
	}

	return keyValue
}
