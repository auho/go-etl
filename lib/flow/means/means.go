package means

type InsertMeans interface {
	GetKeys() []string
	Insert([]string) [][]interface{}
}

type UpdateMeans interface {
	Update([]string) map[string]interface{}
}
