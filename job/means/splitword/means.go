package splitword

import (
	"strings"

	"github.com/auho/go-etl/v2/job/means"
)

var _ means.InsertMeans = (*Means)(nil)

type Means struct {
	SplitWords
	sep string
}

func NewMeans(sep string) *Means {
	s := &Means{}
	s.sep = sep

	return s
}

func (m *Means) GetTitle() string {
	return "SplitWords"
}

func (m *Means) GetKeys() []string {
	return []string{"word"}
}

func (m *Means) Insert(contents []string) []map[string]any {
	items := make([]map[string]any, 0)
	for _, c := range contents {
		results := strings.Split(c, m.sep)
		for _, result := range results {
			items = append(items, map[string]any{"word": result})
		}
	}

	if len(items) <= 0 {
		return nil
	}

	return items
}

func (m *Means) DefaultValues() map[string]any {
	return map[string]any{"word": ""}
}
