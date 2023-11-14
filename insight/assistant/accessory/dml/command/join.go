package command

import (
	"fmt"
	"strings"
)

const symbolFrom = "FROM"
const symbolLeft = "LEFT"

type Join struct {
	symbol  string
	LTable  string
	RTable  string
	LFields []string
	RFields []string
}

func NewJoin(l string, lFields []string, r string, rFields []string) *Join {
	j := &Join{}
	j.symbol = symbolFrom
	j.LTable = l
	j.RTable = r
	j.LFields = lFields
	j.RFields = rFields

	if len(j.LFields) != len(j.RFields) {
		panic(fmt.Sprintf("fields of left[%s] and right[%s] is unequal", strings.Join(j.LFields, ", "), strings.Join(j.RFields, ", ")))
	}

	return j
}

func NewLeftJoin(l string, lFields []string, r string, rFields []string) *Join {
	j := NewJoin(l, lFields, r, rFields)
	j.symbol = symbolLeft

	return j
}

func (j *Join) IsFrom() bool {
	return j.symbol == symbolFrom
}

func (j *Join) IsLeft() bool {
	return j.symbol == symbolLeft
}
