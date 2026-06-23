package search

type FieldSpec interface {
	Keys() []string
	GetDefaultValues() map[string]any
}
