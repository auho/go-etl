package segwordor

import (
	"unicode/utf8"
)

type SegWordsMeans struct {
	SegWords
}

func NewSegWordsMeans() *SegWordsMeans {
	sw := &SegWordsMeans{}
	sw.prepare()

	return sw
}

func (sw *SegWordsMeans) GetTitle() string {
	return "SegWords"
}

func (sw *SegWordsMeans) GetKeys() []string {
	return []string{"word", "flag"}
}

func (sw *SegWordsMeans) Insert(contents []string) []map[string]interface{} {
	results := sw.tag(contents)
	if results == nil {
		return nil
	}

	items := make([]map[string]interface{}, 0, len(results))
	for _, result := range results {
		if utf8.RuneCountInString(result[0]) < 2 || result[1] == "eng" || result[1] == "m" {
			continue
		}

		items = append(items, map[string]interface{}{
			"word": result[0],
			"flag": result[1]},
		)
	}

	return items
}
