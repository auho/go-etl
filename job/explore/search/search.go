package search

type Searcher interface {
	Prepare() error
	GetTitle() string
	GenExport() Exporter
	Do(s []string) Exporter
	Close() error
}
