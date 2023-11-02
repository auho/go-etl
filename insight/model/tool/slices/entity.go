package slices

type valueEntity interface {
	string | int
}

type comparatorEntity interface {
	any | string | int
}
