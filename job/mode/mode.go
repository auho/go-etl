package mode

import (
	"fmt"
	"strings"
	"sync/atomic"

	strings2 "github.com/auho/go-toolkit/farmtools/convert/types/strings"
)

type Operator interface {
	GetTitle() string
	GetFields() []string // source data 里的 key name
	Prepare() error
	Close() error
}

type SingleOperator interface {
	Operator
	Do(map[string]any) map[string]any
}

type InsertOperator interface {
	Operator
	GetKeys() []string             // 处理后的 key name
	DefaultValues() map[string]any // 需要 implement clone important!
	Do(map[string]any) []map[string]any
	State() []string
}

type UpdateOperator interface {
	SingleOperator
}

type TransferOperator interface {
	SingleOperator
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
	s, err := strings2.FromAny(item[key])
	if err != nil {
		panic(fmt.Sprintf("type is not string %T", item[key]))
	}

	return s
}
