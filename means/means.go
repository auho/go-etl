package means

type InsertMeans interface {
	GetName() string
	GetKeys() []string
	Insert([]string) [][]interface{}
}

type UpdateMeans interface {
	Update([]string) map[string]interface{}
}
