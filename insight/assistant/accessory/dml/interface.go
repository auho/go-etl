package dml

type Tabler interface {
	Sql() string
	GetSelectFields() []string
}

type statementor interface {
	Query() string
	InsertQuery(name string) string
	InsertWithFieldsQuery(name string, fields []string) string
	UpdateQuery() string
	DeleteQuery() string
}
