package search

type Searcher interface {
	Prepare() error
	GetTitle() string
	GenExport() Exporter
	Do(s []string) Token
	Close() error
}
