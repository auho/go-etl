package dml

type Tabler interface {
	Sql() string
	GetSelectFields() []string
}
