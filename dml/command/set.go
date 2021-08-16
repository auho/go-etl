package command

import (
	"fmt"
	"strings"
)

type Set struct {
	flag   string
	LTable string
	RTable string
	LKeys  []string
	RKeys  []string
}

func NewSet(lt string, ltKeys []string, rt string, rtKeys []string) *Set {
	s := &Set{}
	s.LTable = lt
	s.RTable = rt
	s.LKeys = ltKeys
	s.RKeys = rtKeys

	if len(s.LKeys) != len(s.RKeys) {
		panic(fmt.Sprintf("keys of left[%s] and right[%s] is unequal", strings.Join(s.LKeys, ", "), strings.Join(s.RKeys, ", ")))
	}

	return s
}

func NewExpressionSet(lt string, ltKeys []string, rt string, rtKeys []string) *Set {
	s := NewSet(lt, ltKeys, rt, rtKeys)
	s.flag = flagExpression

	return s
}

func (s *Set) IsExpression() bool {
	return s.flag == flagExpression
}
