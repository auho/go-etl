package segword

import (
	"github.com/auho/go-etl/v2/job/explore/search"
)

var _ search.Searcher = (*SegWords)(nil)

type SegWords struct {
	seg    *Seg
	export *Export
}

func NewSegWords(export *Export) *SegWords {
	return &SegWords{export: export}
}

func (sg *SegWords) GetTitle() string {
	return "Seg"
}

func (sg *SegWords) GenExport() search.Exporter {
	return sg.export
}

func (sg *SegWords) Prepare() error {
	sg.seg = NewSeg()

	return nil
}

func (sg *SegWords) Do(contents []string) search.Token {
	results := sg.seg.tag(contents)

	return sg.export.ToToken(results)
}

func (sg *SegWords) Close() error {
	return sg.seg.Close()
}
