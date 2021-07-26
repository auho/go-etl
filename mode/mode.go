package mode

type Mode interface {
}

type VoidMode interface {
	Do(map[string]interface{})
}

type InsertMode interface {
	GetKeys() []string
	GetFields() []string
	Do(map[string]interface{}) [][]interface{}
}

type UpdateMode interface {
	GetFields() []string
	Do(map[string]interface{}) map[string]interface{}
}
