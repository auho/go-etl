package search

type Exporter interface {
	IsOk() bool
	GetKeys() []string
	GetDefaultValues() map[string]any
	ToTokenize() []map[string]any
}
