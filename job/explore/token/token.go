package token

type Tokenizer interface {
	GetOk() bool
	ToModer(m string) Moder
}

type Moder interface {
	DefaultValues() map[string]any
	ToTokenize() []map[string]any
}

type Tokenize struct {
}
