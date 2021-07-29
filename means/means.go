package means

type InsertMeans interface {
	GetKeys() []string
	Insert([]string) [][]interface{}
	Close()
}

type UpdateMeans interface {
	Update([]string) map[string]interface{}
	Close()
}
