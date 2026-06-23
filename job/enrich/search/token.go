package search

type Token struct {
	ok            bool
	tokenizerFunc func() []map[string]any
}

func (t *Token) SetOK() {
	t.ok = true
}

func (t *Token) SetTokenizerFunc(fn func() []map[string]any) {
	t.tokenizerFunc = fn
}

func (t *Token) IsOK() bool {
	return t.ok
}

func (t *Token) ToToken() []map[string]any {
	if t.tokenizerFunc == nil {
		return nil
	}

	return t.tokenizerFunc()
}
