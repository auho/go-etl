package filter

// Predicate decides whether an item satisfies a condition.
//
// Prepare initializes resources the predicate may hold (for example a Collector
// that needs warming up) and returns any error so it can propagate up the
// pipeline instead of being swallowed or causing a panic.
type Predicate interface {
	Prepare() error
	Match(map[string]any) (bool, error)
}

// Func is a function adapter that satisfies Predicate, the Predicate analogue
// of http.HandlerFunc. Use it to wrap a plain function when no Prepare step is
// needed.
type Func func(map[string]any) (bool, error)

func (f Func) Prepare() error { return nil }

func (f Func) Match(m map[string]any) (bool, error) { return f(m) }
