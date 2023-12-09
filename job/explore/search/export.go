package search

type Exporter interface {
	IsOk() bool
	DefaultValues() map[string]any
	ToTokenize() []map[string]any
}
