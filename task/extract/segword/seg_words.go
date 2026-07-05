package segword

import (
	"unicode/utf8"

	"github.com/auho/go-etl/v3/task/extract"
)

var DefaultFilterFunc = func(result result) bool {
	return utf8.RuneCountInString(result.token) < 2 || result.flag == "eng" || result.flag == "m"
}

var _ extract.Extractor = (*SegWords)(nil)

type SegWords struct {
	seg        *Seg
	format     format
	toMaps     func(results, format) []map[string]any
	filterFunc func(result) bool
	keys       []string
	defaults   map[string]any
}

func NewSegWordsAll() *SegWords {
	fm := defaultFormat
	return &SegWords{
		format:     fm,
		toMaps:     func(r results, f format) []map[string]any { return r.toAll(f) },
		filterFunc: DefaultFilterFunc,
		keys:       []string{fm.tokenName, fm.flagName},
		defaults:   map[string]any{fm.tokenName: "", fm.flagName: ""},
	}
}

func NewSegWordsLine() *SegWords {
	fm := defaultFormat
	return &SegWords{
		format:     fm,
		toMaps:     func(r results, f format) []map[string]any { return r.toLine(f) },
		filterFunc: DefaultFilterFunc,
		keys:       []string{fm.tokenName},
		defaults:   map[string]any{fm.tokenName: ""},
	}
}

func (sg *SegWords) Title() string {
	return "Seg"
}

func (sg *SegWords) Prepare() error {
	sg.seg = NewSeg()
	sg.format.check()
	return nil
}

func (sg *SegWords) Keys() []string {
	return sg.keys
}

func (sg *SegWords) DefaultValues() map[string]any {
	return sg.defaults
}

func (sg *SegWords) Extract(contents []string) extract.Result {
	all := sg.seg.tag(contents)

	var filtered results
	for _, r := range all {
		if !sg.filterFunc(r) {
			filtered = append(filtered, r)
		}
	}

	if len(filtered) == 0 {
		return extract.Result{}
	}
	return extract.NewResult(true, sg.toMaps(filtered, sg.format))
}

func (sg *SegWords) Close() error {
	return sg.seg.Close()
}
