package segword

import (
	"unicode/utf8"

	"github.com/auho/go-etl/v2/job/means"
)

var _ means.InsertMeans = (*Means)(nil)

type Means struct {
	SegWords
}

func NewMeans() *Means {
	sw := &Means{}
	sw.prepare()

	return sw
}

func (m *Means) GetTitle() string {
	return "SegWords"
}

func (m *Means) GetKeys() []string {
	return []string{"word", "flag"}
}

func (m *Means) Insert(contents []string) []map[string]any {
	results := m.tag(contents)
	if results == nil {
		return nil
	}

	items := make([]map[string]any, 0, len(results))
	for _, result := range results {
		if utf8.RuneCountInString(result[0]) < 2 || result[1] == "eng" || result[1] == "m" {
			continue
		}

		items = append(items, map[string]any{
			"word": result[0],
			"flag": result[1]},
		)
	}

	return items
}

func (m *Means) DefaultValues() map[string]any {
	return map[string]any{"word": "", "flag": ""}
}
