package search

type Searcher interface {
	GetTitle() string
	Prepare() error
	GenExport() FieldSpec
	Do(s []string) Token
	Close() error
}
