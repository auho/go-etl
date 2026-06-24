package search

type FieldSpec interface {
	Keys() []string
	DefaultValues() map[string]any
}
