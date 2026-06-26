package extract

type Extractor interface {
	Title() string
	Prepare() error
	NewExport() FieldSpec
	Search(contents []string) Result
	Close() error
}
