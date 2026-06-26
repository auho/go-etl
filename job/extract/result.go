package extract

type Result struct {
	ok          bool
	resultsFunc func() []map[string]any
}

func (t *Result) SetOK() {
	t.ok = true
}

func (t *Result) SetResultsFunc(fn func() []map[string]any) {
	t.resultsFunc = fn
}

func (t *Result) IsOK() bool {
	return t.ok
}

func (t *Result) Rows() []map[string]any {
	if t.resultsFunc == nil {
		return nil
	}

	return t.resultsFunc()
}
