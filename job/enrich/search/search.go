package search

type Searcher interface {
	Title() string
	Prepare() error
	GenExport() FieldSpec
	Do(s []string) Token
	Close() error
}
