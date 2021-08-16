package command

type Query interface {
	FieldsForInsert() []string
	Sql() string
}

type TableJoinCommander interface {
	BuildFieldsForInsert() []string
	Query() string
}

type TableCommander interface {
	BuildFieldsForInsert() []string
	SetTable(string, string)
	SetSelect(*Entries)
	Select() string
	BuildSelect() []string
	SetFrom(*Join)
	From() string
	BuildFrom() []string
	SetWhere(string)
	Where() string
	BuildWhere() []string
	SetGroupBy(*Entries)
	GroupBy() string
	BuildGroupBy() []string
	SetOrderBy(*Entries)
	OrderBy() string
	BuildOrderBy() []string
	SetLimit([]int)
	Limit() string
	Query() string
}

type InsertCommander interface {
	SetTable(string)
	SetFields([]string)
	SetQuery(Query)
	Insert() string
	InsertWithFields() string
}
