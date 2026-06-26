package filter

var _ Spec = (*AND)(nil)
var _ Spec = (*OR)(nil)

type Predicate func(map[string]any) bool

type Expression []Predicate

type AND Expression

func NewAND(ops ...Predicate) AND {
	a := AND{}
	a = append(a, ops...)

	return a
}

func (a AND) OK(item map[string]any) bool {
	for _, op := range a {
		if !op(item) {
			return false
		}
	}

	return true
}

func (a AND) ToPredicate() Predicate {
	return func(m map[string]any) bool {
		return a.OK(m)
	}
}

type OR Expression

func NewOR(ops ...Predicate) OR {
	o := OR{}
	o = append(o, ops...)

	return o
}

func (o OR) OK(item map[string]any) bool {
	for _, op := range o {
		if op(item) {
			return true
		}
	}

	return false
}

func (o OR) ToPredicate() Predicate {
	return func(m map[string]any) bool {
		return o.OK(m)
	}
}
