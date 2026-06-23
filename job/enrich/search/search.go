package search

type Searcher interface {
	GetTitle() string
	Prepare() error
	GenExport() Exporter
	Do(s []string) Token
	Close() error
}
