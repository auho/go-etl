package command

type Query interface {
	BuildFieldsForInsert() []string
	Query() string
}

type TableJoinCommander interface {
	SetCommands([]TableCommander)
	SetLimit(l []int)
	Query() string
	InsertQuery(name string) string
	InsertWithFieldsQuery(name string, fields []string) string
	UpdateQuery() string
	DeleteQuery() string
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
	SetSet([]*Set)
	Set() string
	BuildSet() []string
	Query() string
	InsertQuery(name string) string
	InsertWithFieldsQuery(name string, fields []string) string
	UpdateQuery() string
	DeleteQuery() string
}
