package command

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v3/tool/slicex"
)

type Set struct {
	flag    string   // = 后的参数值类型标识 key ｜ expression ｜ value
	LTable  string   // = 前的 table name
	RTable  string   // = 后的 table name
	LFields []string // = 前的 field name
	RValues []any    // = 后的 field name | expression | value
}

func newSet(l string, fields []string, r string, values []any, flag string) *Set {
	s := &Set{}
	s.LTable = l
	s.RTable = r
	s.LFields = fields
	s.RValues = values
	s.flag = flag

	if len(s.LFields) != len(s.RValues) {
		var _rvs []string
		for _, _rv := range s.RValues {
			_rvs = append(_rvs, fmt.Sprintf("%v", _rv))
		}

		panic(fmt.Sprintf("fields of left[%s] and right[%s] is unequal", strings.Join(s.LFields, ", "), strings.Join(_rvs, ", ")))
	}

	return s
}

// NewSetField
// 传入 field name => field name
func NewSetField(l string, lFields []string, r string, rFields []string) *Set {
	return newSet(l, lFields, r, slicex.SliceToAny(rFields), "")
}

// NewSetExpression
// 传入 field name => expression
func NewSetExpression(l string, lFields []string, r string, expression []string) *Set {
	return newSet(l, lFields, r, slicex.SliceToAny(expression), FlagExpression)
}

// NewSetValue
// 传入 field name => value
func NewSetValue(l string, lFields []string, r string, values []any) *Set {
	return newSet(l, lFields, r, values, FlagValue)
}

func (s *Set) IsExpression() bool {
	return s.flag == FlagExpression
}

func (s *Set) IsValue() bool {
	return s.flag == FlagValue
}
