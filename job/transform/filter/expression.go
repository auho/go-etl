package filter

import "fmt"

var _ Spec = (*And)(nil)
var _ Spec = (*Or)(nil)

type Expression = []Predicate

type And Expression

func NewAnd(ps ...Predicate) And {
	a := And{}
	a = append(a, ps...)

	return a
}

func (a And) OK(m map[string]any) (bool, error) {
	for _, op := range a {
		ok, err := op(m)
		if err != nil {
			return false, fmt.Errorf("And.OK: %w", err)
		}
		if !ok {
			return false, nil
		}
	}

	return true, nil
}

func (a And) ToPredicate() Predicate {
	return func(m map[string]any) (bool, error) {
		return a.OK(m)
	}
}

type Or Expression

func NewOr(ps ...Predicate) Or {
	o := Or{}
	o = append(o, ps...)

	return o
}

func (o Or) OK(m map[string]any) (bool, error) {
	for _, op := range o {
		ok, err := op(m)
		if err != nil {
			return false, fmt.Errorf("Or.OK: %w", err)
		}
		if ok {
			return true, nil
		}
	}

	return false, nil
}

func (o Or) ToPredicate() Predicate {
	return func(m map[string]any) (bool, error) {
		return o.OK(m)
	}
}
