package command

import (
	"fmt"
	"strings"
)

const symbolFrom = "FROM"
const symbolLeft = "LEFT"

type Join struct {
	symbol string
	LTable string
	RTable string
	LKeys  []string
	RKeys  []string
}

func NewJoin(lt string, ltKeys []string, rt string, rtKeys []string) *Join {
	j := &Join{}
	j.symbol = symbolFrom
	j.LTable = lt
	j.RTable = rt
	j.LKeys = ltKeys
	j.RKeys = rtKeys

	if len(j.LKeys) != len(j.RKeys) {
		panic(fmt.Sprintf("keys of left[%s] and right[%s] is unequal", strings.Join(j.LKeys, ", "), strings.Join(j.RKeys, ", ")))
	}

	return j
}

func NewLeftJoin(lt string, ltKeys []string, rt string, rtKeys []string) *Join {
	j := NewJoin(lt, ltKeys, rt, rtKeys)
	j.symbol = symbolLeft

	return j
}

func (j *Join) IsFrom() bool {
	return j.symbol == symbolFrom
}

func (j *Join) IsLeft() bool {
	return j.symbol == symbolLeft
}
