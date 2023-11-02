package slices

type Any interface {
	string | int
}

type comparator interface {
	any | string | int
}
