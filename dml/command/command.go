package command

type tableCommander interface {
	Select(map[string]string) string
	BuildSelect(map[string]string) []string
	From(*Join) string
	BuildFrom(j *Join) []string
	Where(string) string
	BuildWhere(s string) []string
	GroupBy(map[string]string) string
	BuildGroupBy(map[string]string) []string
	OrderBy(map[string]string) string
	BuildOrderBy(map[string]string) []string
	Limit([]int) string
}

type TableJoinCommander interface {
	SelectToString([]string) string
	FromToString([]string) string
	WhereToString([]string) string
	GroupByToString([]string) string
	OrderByToString([]string) string
	LimitToString([]int) string
}
