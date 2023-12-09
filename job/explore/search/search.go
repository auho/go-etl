package search

type Searcher interface {
	GetTitle() string
	InitialExporter() Exporter
	Do(s []string) Exporter
}
