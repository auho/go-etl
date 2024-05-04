package search

type Token struct {
	ok            bool
	tokenizerFunc func() []map[string]any
}

func (t *Token) SetOk() {
	t.ok = true
}

func (t *Token) SetTokenizerFunc(fn func() []map[string]any) {
	t.tokenizerFunc = fn
}

func (t *Token) IsOk() bool {
	return t.ok
}

func (t *Token) ToToken() []map[string]any {
	if t.tokenizerFunc == nil {
		return nil
	}

	return t.tokenizerFunc()
}
