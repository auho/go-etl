package filter

import "fmt"

var _ Spec = (*AND)(nil)
var _ Spec = (*OR)(nil)

type Predicate func(map[string]any) (bool, error)

type Expression []Predicate

type AND Expression

func NewAND(ops ...Predicate) AND {
	a := AND{}
	a = append(a, ops...)

	return a
}

func (a AND) OK(item map[string]any) (bool, error) {
	for _, op := range a {
		ok, err := op(item)
		if err != nil {
			return false, fmt.Errorf("AND.OK: %w", err)
		}
		if !ok {
			return false, nil
		}
	}

	return true, nil
}

func (a AND) ToPredicate() Predicate {
	return func(m map[string]any) (bool, error) {
		return a.OK(m)
	}
}

type OR Expression

func NewOR(ops ...Predicate) OR {
	o := OR{}
	o = append(o, ops...)

	return o
}

func (o OR) OK(item map[string]any) (bool, error) {
	for _, op := range o {
		ok, err := op(item)
		if err != nil {
			return false, fmt.Errorf("OR.OK: %w", err)
		}
		if ok {
			return true, nil
		}
	}

	return false, nil
}

func (o OR) ToPredicate() Predicate {
	return func(m map[string]any) (bool, error) {
		return o.OK(m)
	}
}
