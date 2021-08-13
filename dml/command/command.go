package command

type SqlCommander interface {
	Sql() string
}

type DriverCommander interface {
	SelectToString([]string) string
	FromToString([]string) string
	WhereToString([]string) string
	GroupByToString([]string) string
	OrderByToString([]string) string
	LimitToString([]int) string
}

type TableCommander interface {
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
}

type InsertCommander interface {
	SetTable(string)
	SetFields([]string)
	Insert() string
}
