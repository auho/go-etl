package extract

// Extractor transforms content strings into structured rows.
// Keys and DefaultValues are prepared during Prepare so callers can
// retrieve them without re-computing.
type Extractor interface {
	Prepare() error
	Title() string
	Keys() []string
	DefaultValues() map[string]any
	Extract(contents []string) Result
	Close() error
}
