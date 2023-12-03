package maps

type keyEntity interface {
	string | int
}

type valueEntity interface {
	string | int
}

type comparatorEntity interface {
	any | string | int
}
