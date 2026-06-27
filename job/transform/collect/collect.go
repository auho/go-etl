package collect

import (
	"fmt"

	"github.com/auho/go-etl/v3/job/extract"
	"github.com/auho/go-toolkit/v2/farmtools/convert/types/strings"
)

type Collector interface {
	Title() string
	Keys() []string // for source select data row
	Search(item map[string]any, search extract.Extractor) (extract.Result, error)
}

type Collect struct{}

func (c *Collect) GetKeyContent(key string, item map[string]any) (string, error) {
	s, err := strings.FromAny(item[key])
	if err != nil {
		return "", fmt.Errorf("GetKeyContent[%s]: %w", key, err)
	}

	return s, nil
}

func (c *Collect) GetKeysContent(keys []string, item map[string]any) ([]string, error) {
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
