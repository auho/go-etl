package command

type Driver interface {
	Query() string
	FieldsForInsert() string
}

type TableJoinCommander interface {
	SelectToString([]string) string
	FromToString([]string) string
	WhereToString([]string) string
	GroupByToString([]string) string
	OrderByToString([]string) string
	LimitToString([]int) string
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
	Insert() string
}
