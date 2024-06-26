package collect

import (
	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-toolkit/farmtools/convert/types/strings"
)

type Collector interface {
	GetTitle() string
	GetKeys() []string // for source select data row
	Do(item map[string]any, search search.Searcher) search.Token
}

type Collect struct{}

func (c *Collect) GetKeyContent(key string, item map[string]any) string {
	s, err := strings.FromAny(item[key])
	if err != nil {
		panic(err)
	}

	return s
}

func (c *Collect) GetKeysContent(keys []string, item map[string]any) []string {
	contents := make([]string, 0)
	for _, key := range keys {
		keyValue, err := strings.FromAny(item[key])
		if err != nil {
			panic(err)
		}

		contents = append(contents, keyValue)
	}

	return contents
}
