package extract

// Result holds the outcome of an extraction.
// ok reports whether the extraction found anything; rows holds the
// structured output. Both are set at construction time — there is no
// lazy evaluation.
type Result struct {
	ok   bool
	rows []map[string]any
}

// NewResult builds a Result from its components.
func NewResult(ok bool, rows []map[string]any) Result {
	return Result{ok: ok, rows: rows}
}

// Get returns the ok flag and rows in one call.
func (r Result) Get() (bool, []map[string]any) {
	return r.ok, r.rows
}

// IsOK reports whether the extraction found anything.
func (r Result) IsOK() bool {
	return r.ok
}
