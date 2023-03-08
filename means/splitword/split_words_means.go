package splitword

import (
	"strings"
)

type SplitWordsMeans struct {
	sep string
}

func NewSplitWordsMeans(sep string) *SplitWordsMeans {
	s := &SplitWordsMeans{}
	s.sep = sep

	return s
}

func (s *SplitWordsMeans) GetTitle() string {
	return "SplitWords"
}

func (s *SplitWordsMeans) GetKeys() []string {
	return []string{"word"}
}

func (s *SplitWordsMeans) Insert(contents []string) []map[string]any {
	items := make([]map[string]any, 0)
	for _, c := range contents {
		results := strings.Split(c, s.sep)
		for _, result := range results {
			items = append(items, map[string]any{"word": result})
		}
	}

	if len(items) <= 0 {
		return nil
	}

	return items
}

func (s *SplitWordsMeans) Close() {}
