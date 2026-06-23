package search

type FieldSpec interface {
	GetKeys() []string
	GetDefaultValues() map[string]any
}
