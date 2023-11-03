package tag

// #TODO 重命名（key keyword  label=>tag） SQL use orm or interface

// Result
// result 匹配结果
type Result struct {
	Key        string            // keyword
	Num        int64             // matched num
	Texts      map[string]int64  // matched text map[matched text]num
	Tags       map[string]string // tags map[tag name]tag
	IsKeyMerge bool
}

func NewResult() *Result {
	m := &Result{}
	m.Tags = make(map[string]string)
	m.Texts = make(map[string]int64)

	return m
}

// LabelResult
// label result
type LabelResult struct {
	Identity string
	Labels   map[string]string         // tags map[tag name]tag
	Match    map[string]map[string]int // keyword and match text map[keyword]map[matched text]num
	KeyNum   int64                     // keyword num
	TextNum  int64                     // all text word num
}

func NewLabelResult() *LabelResult {
	l := &LabelResult{}
	l.Labels = make(map[string]string)
	l.Match = make(map[string]map[string]int)

	return l
}

// Results
// results
type Results []*Result

func (r Results) Len() int {
	return len(r)
}

func (r Results) Less(i, j int) bool {
	return r[i].Num > r[j].Num
}

func (r Results) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// LabelResults
// label results
type LabelResults []*LabelResult

func (l LabelResults) Len() int {
	return len(l)
}

func (l LabelResults) Less(i, j int) bool {
	return l[i].TextNum > l[j].TextNum
}

func (l LabelResults) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
