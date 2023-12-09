package search

type Searcher interface {
	GetTitle() string
	GetExport() Exporter
	Do(s []string) Exporter
}
