package means

type InsertMeans interface {
	GetTitle() string
	GetKeys() []string
	Insert([]string) []map[string]any
	Prepare() error
	Close() error
}

type UpdateMeans interface {
	GetTitle() string
	Update([]string) map[string]any
	Prepare() error
	Close() error
}
