package filter

type Predicate func(map[string]any) (bool, error)

type Spec interface {
	OK(map[string]any) (bool, error)
	ToPredicate() Predicate
}
