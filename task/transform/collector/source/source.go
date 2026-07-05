package source

// Source provides ordered content values from a data source.
type Source interface {
	Title() string
	Keys() []string
	Prepare() error
	Contents(item map[string]any) (keys []string, keysValue map[string]string, err error)
}
