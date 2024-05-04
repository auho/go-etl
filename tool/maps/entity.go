package maps

type K interface {
	string | int
}

type V interface {
	any | string | int
}

type comparatorV interface {
	any | string | int
}
