package collect

import (
	"fmt"

	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-toolkit/v2/farmtools/convert/types/strings"
)

type Collector interface {
	Title() string
	SourceKeys() []string // for source select data row
	Extract(item map[string]any, e extract.Extractor) (extract.Result, error)
}

type ContentReader struct{}

func (c *ContentReader) Content(key string, item map[string]any) (string, error) {
	s, err := strings.FromAny(item[key])
	if err != nil {
		return "", fmt.Errorf("Content[%s]: %w", key, err)
	}

	return s, nil
}

func (c *ContentReader) Contents(keys []string, item map[string]any) ([]string, error) {
	contents := make([]string, 0)
	for _, key := range keys {
		keyValue, err := strings.FromAny(item[key])
		if err != nil {
			return nil, fmt.Errorf("FromAny[%s]: %w", key, err)
		}

		contents = append(contents, keyValue)
	}

	return contents, nil
}
