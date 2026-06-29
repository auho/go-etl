package filter

import "fmt"

var _ Predicate = And{}
var _ Predicate = Or{}

// Expression is an ordered list of predicates used by the logical combinators
// And and Or.
type Expression = []Predicate

// And is a predicate that is true only when every operand matches.
type And Expression

func NewAnd(ps ...Predicate) And {
	a := And{}
	a = append(a, ps...)

	return a
}

func (a And) Prepare() error {
	for _, op := range a {
		if err := op.Prepare(); err != nil {
			return err
		}
	}

	return nil
}

func (a And) Match(m map[string]any) (bool, error) {
	for _, op := range a {
		ok, err := op.Match(m)
		if err != nil {
			return false, fmt.Errorf("And.Match: %w", err)
		}
		if !ok {
			return false, nil
		}
	}

	return true, nil
}

// Or is a predicate that is true when any operand matches.
type Or Expression

func NewOr(ps ...Predicate) Or {
	o := Or{}
	o = append(o, ps...)

	return o
}

func (o Or) Prepare() error {
	for _, op := range o {
		if err := op.Prepare(); err != nil {
			return err
		}
	}

	return nil
}

func (o Or) Match(m map[string]any) (bool, error) {
	for _, op := range o {
		ok, err := op.Match(m)
		if err != nil {
			return false, fmt.Errorf("Or.Match: %w", err)
		}
		if ok {
			return true, nil
		}
	}

	return false, nil
}
