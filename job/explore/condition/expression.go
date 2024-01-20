package condition

var _ Conditioner = (*AND)(nil)
var _ Conditioner = (*OR)(nil)

type Operation func(map[string]any) bool

type Expression []Operation

type AND Expression

func NewAND(ops ...Operation) AND {
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

func (a AND) ToOperation() Operation {
	return func(m map[string]any) bool {
		return a.OK(m)
	}
}

type OR Expression

func NewOR(ops ...Operation) OR {
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

func (o OR) ToOperation() Operation {
	return func(m map[string]any) bool {
		return o.OK(m)
	}
}
