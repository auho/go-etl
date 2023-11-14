package command

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/tool/slices"
)

type Set struct {
	flag    string   // = 后的参数值类型标识 key ｜ expression ｜ value
	LTable  string   // = 前的 table name
	RTable  string   // = 后的 table name
	LKeys   []string // = 前的 field name
	RValues []any    // = 后的 field name | expression | value
}

func newSet(lt string, ltKeys []string, rt string, rtValue []any, flag string) *Set {
	s := &Set{}
	s.LTable = lt
	s.RTable = rt
	s.LKeys = ltKeys
	s.RValues = rtValue
	s.flag = flag

	if len(s.LKeys) != len(s.RValues) {
		var _rvs []string
		for _, _rv := range s.RValues {
			_rvs = append(_rvs, fmt.Sprintf("%v", _rv))
		}

		panic(fmt.Sprintf("keys of left[%s] and right[%s] is unequal", strings.Join(s.LKeys, ", "), strings.Join(_rvs, ", ")))
	}

	return s
}

// NewSetField
// 传入 field name => field name
func NewSetField(l string, lFields []string, r string, rFields []string) *Set {
	return newSet(l, lFields, r, slices.SliceToAny(rFields), "")
}

// NewSetExpression
// 传入 field name => expression
func NewSetExpression(l string, lFields []string, r string, expression []string) *Set {
	return newSet(l, lFields, r, slices.SliceToAny(expression), flagExpression)
}

// NewSetValue
// 传入 field name => value
func NewSetValue(l string, lFields []string, r string, values []any) *Set {
	return newSet(l, lFields, r, values, flagValue)
}

func (s *Set) IsExpression() bool {
	return s.flag == flagExpression
}

func (s *Set) IsValue() bool {
	return s.flag == flagValue
}
