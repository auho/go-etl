package command

type Query interface {
	BuildFieldsForInsert() []string
	Query() string
}

type TableJoinCommander interface {
	SetCommands([]TableCommander)
	SetLimit(l []int)
	Query() string
	Insert(name string) string
	InsertWithFields(name string, fields []string) string
	Delete() string
}

type TableCommander interface {
	BuildFieldsForInsert() []string
	Name() string
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
	Insert(name string) string
	InsertWithFields(name string, fields []string) string
	Delete() string
}

type InsertCommander interface {
	SetTable(string)
	SetFields([]string)
	SetQuery(Query)
	Insert() string
	InsertWithFields() string
}
