package transform

import (
	"fmt"
	"strings"
	"sync/atomic"
)

type base struct {
	keys   []string // 要被处理的 key name
	total  int64
	amount int64
}

func (m *base) addTotal(num int64) {
	atomic.AddInt64(&m.total, num)
}

func (m *base) addAmount(num int64) {
	atomic.AddInt64(&m.amount, num)
}

func (m *base) genCounter() string {
	return fmt.Sprintf("total: %d; amount: %d", atomic.LoadInt64(&m.total), atomic.LoadInt64(&m.amount))
}

func (m *base) genTitle(name string, desc string) string {
	return fmt.Sprintf("%s %s{%s}", name, "keys["+strings.Join(m.keys, ", ")+"]", desc)
}
