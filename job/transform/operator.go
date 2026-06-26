package transform

import (
	"fmt"
	"strings"
	"sync/atomic"

	strings2 "github.com/auho/go-toolkit/v2/farmtools/convert/types/strings"
)

type Operator interface {
	Title() string
	GetFields() []string // source data 里的 key name
	Prepare() error
	Close() error
}

type SingleOperator interface {
	Operator
	Apply(map[string]any) map[string]any
}

type InsertOperator interface {
	Operator
	Keys() []string                // 处理后的 key name
	DefaultValues() map[string]any // 需要 implement clone important!
	Apply(map[string]any) []map[string]any
	State() []string
}

type UpdateOperator interface {
	SingleOperator
}

type TransferOperator interface {
	SingleOperator
}

type operator struct {
	keys   []string // 要被处理的 key name
	total  int64
	amount int64
}

func (o *operator) AddTotal(num int64) {
	atomic.AddInt64(&o.total, num)
}

func (o *operator) AddAmount(num int64) {
	atomic.AddInt64(&o.amount, num)
}

func (o *operator) GenCounter() string {
	return fmt.Sprintf("total: %d; amount: %d", o.total, o.amount)
}

func (o *operator) GenTitle(name string, desc string) string {
	return fmt.Sprintf("%s %s{%s}", name, "keys["+strings.Join(o.keys, ", ")+"]", desc)
}

func (o *operator) GetKeyContent(key string, item map[string]any) string {
	return o.KeyValueToString(key, item)
}

func (o *operator) GetKeysContent(keys []string, item map[string]any) []string {
	contents := make([]string, 0)
	for _, key := range keys {
		keyValue := o.KeyValueToString(key, item)

		contents = append(contents, keyValue)
	}

	return contents
}

func (o *operator) KeyValueToString(key string, item map[string]any) string {
	s, err := strings2.FromAny(item[key])
	if err != nil {
		panic(fmt.Sprintf("type is not string %T", item[key]))
	}

	return s
}
