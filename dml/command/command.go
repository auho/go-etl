package command

type commander interface {
	Select(map[string]string) string
	From(*Join) string
	Where(string) string
	GroupBy(map[string]string) string
	OrderBy(map[string]string) string
	Limit([]int) string
	TableName() string
}
