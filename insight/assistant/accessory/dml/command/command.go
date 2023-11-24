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
	DriverName() string
	BuildFieldsForInsert() []string
	Name() string
	SetTable(string, string)
	SetSelect(*Entities)
	Select() string
	BuildSelect() []string
	SetFrom(*Join)
	From() string
	BuildFrom() []string
	SetWhere(string)
	Where() string
	BuildWhere() []string
	SetGroupBy(*Entities)
	GroupBy() string
	BuildGroupBy() []string
	SetOrderBy(*Entities)
	OrderBy() string
	BuildOrderBy() []string
	SetLimit([]int)
	Limit() string
	SetHaving(string)
	Having() string
	SetSet([]*Set)
	Set() string
	BuildSet() []string
	Query() string
	InsertQuery(name string) string
	InsertWithFieldsQuery(name string, fields []string) string
	UpdateQuery() string
	DeleteQuery() string
}
