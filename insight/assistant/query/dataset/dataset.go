package dataset

// Dataset represents the raw data set before processing.
type Dataset struct {
	Name   string   // dataset name
	Keys   []string // item key name
	Titles []string // dataset item data title
	Sets   []Subset
}

// Result represents the processed data after merging.
type Result struct {
	Names      []string           // data names, preserves name order
	Rows       map[string][][]any // map[name]rows
	RowsAmount map[string]int     // rows count (excluding title)
	Amount     int                // total rows count (excluding title)
}

// NewResult creates a Result with initialized maps.
func NewResult() *Result {
	return &Result{
		Rows:       make(map[string][][]any),
		RowsAmount: make(map[string]int),
	}
}

func (r *Result) addRows(name string, rows [][]any) {
	r.Names = append(r.Names, name)
	r.RowsAmount[name] = len(rows)
	r.Amount += r.RowsAmount[name]
	r.Rows[name] = rows
}

func (r *Result) addRowsWithTitles(name string, titles []any, rows [][]any) {
	r.Names = append(r.Names, name)
	r.RowsAmount[name] = len(rows)
	r.Amount += r.RowsAmount[name]
	r.Rows[name] = append([][]any{titles}, rows...)
}
