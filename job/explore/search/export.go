package search

type Exporter interface {
	IsOk() bool
	GetKeys() []string
	DefaultValues() map[string]any
	ToTokenize() []map[string]any
}
