package segword

import (
	"unicode/utf8"

	"github.com/auho/go-etl/v3/job/extract"
)

var DefaultFilterFunc = func(result Result) bool {
	return utf8.RuneCountInString(result.Token) < 2 || result.Flag == "eng" || result.Flag == "m"
}

var _ extract.Extractor = (*SegWords)(nil)

type SegWords struct {
	seg        *Seg
	export     *extract.Exporter[Results]
	filterFunc func(Result) bool
}

func NewDefault() *SegWords {
	return NewSegWords(NewExportAll())
}

func NewSegWords(export *extract.Exporter[Results]) *SegWords {
	return &SegWords{export: export, filterFunc: DefaultFilterFunc}
}

func (sg *SegWords) Title() string {
	return "Seg"
}

func (sg *SegWords) NewExport() extract.FieldSpec {
	return sg.export
}

func (sg *SegWords) Prepare() error {
	sg.seg = NewSeg()

	return nil
}

func (sg *SegWords) Search(contents []string) extract.Result {
	results := sg.seg.tag(contents)

	var filtered Results
	for _, r := range results {
		if !sg.filterFunc(r) {
			filtered = append(filtered, r)
		}
	}

	return sg.export.ToToken(filtered, len(filtered) > 0)
}

func (sg *SegWords) Close() error {
	return sg.seg.Close()
}
