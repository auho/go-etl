package search

type Exporter interface {
	GetKeys() []string
	GetDefaultValues() map[string]any
}

type Token struct {
	Ok        bool
	Tokenizer func() []map[string]any
}

func (t *Token) IsOk() bool {
	return t.Ok
}

func (t *Token) ToToken() []map[string]any {
	return t.Tokenizer()
}
