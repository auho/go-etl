package search

type Exporter interface {
	GetKeys() []string
	GetDefaultValues() map[string]any
}
