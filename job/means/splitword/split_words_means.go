package splitword

import (
	"strings"

	"github.com/auho/go-etl/v2/job/means"
)

var _ means.InsertMeans = (*SplitWordsMeans)(nil)

type SplitWordsMeans struct {
	SplitWords
	sep string
}

func NewSplitWordsMeans(sep string) *SplitWordsMeans {
	s := &SplitWordsMeans{}
	s.sep = sep

	return s
}

func (sw *SplitWordsMeans) GetTitle() string {
	return "SplitWords"
}

func (sw *SplitWordsMeans) GetKeys() []string {
	return []string{"word"}
}

func (sw *SplitWordsMeans) Insert(contents []string) []map[string]any {
	items := make([]map[string]any, 0)
	for _, c := range contents {
		results := strings.Split(c, sw.sep)
		for _, result := range results {
			items = append(items, map[string]any{"word": result})
		}
	}

	if len(items) <= 0 {
		return nil
	}

	return items
}

func (sw *SplitWordsMeans) DefaultValues() map[string]any {
	return map[string]any{"word": ""}
}
