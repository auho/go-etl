package token

type Tokenizer interface {
	GetOk() bool
	ToExport(m string) Exporter
}

type Exporter interface {
	Name() string
	DefaultValues() map[string]any
	ToTokenize() []map[string]any
}

type Tokenize struct {
}
