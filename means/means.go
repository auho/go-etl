package means

type InsertMeans interface {
	GetTitle() string
	GetKeys() []string
	Insert([]string) []map[string]interface{}
	Close()
}

type UpdateMeans interface {
	GetTitle() string
	Update([]string) map[string]interface{}
	Close()
}
