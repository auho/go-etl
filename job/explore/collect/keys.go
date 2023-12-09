package collect

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/job/explore/search"
)

var _ Collector = (*Keys)(nil)

const (
	keysWayAll = iota
	keysWayInOrder
)

// Keys
// collect from keys
type Keys struct {
	Collect

	keys []string
	way  int
}

// NewKeys
// all keys
func NewKeys(keys []string) *Keys {
	return newKeys(keys, keysWayAll)
}

// NewKeysInOrder
// one by one, if matched return
func NewKeysInOrder(keys []string) *Keys {
	return newKeys(keys, keysWayInOrder)
}

func newKeys(keys []string, way int) *Keys {
	return &Keys{
		keys: keys,
		way:  way,
	}
}

func (f *Keys) GetTitle() string {
	return fmt.Sprintf("keys{%s}", strings.Join(f.keys, ","))
}

func (f *Keys) GetKeys() []string {
	return f.keys
}

func (f *Keys) Do(item map[string]any, searcher search.Searcher) search.Exporter {
	if f.IsAll() {
		return f.doAll(item, searcher)
	} else if f.IsInOrder() {
		return f.doInOrder(item, searcher)
	} else {
		panic("way unknown")
	}
}

func (f *Keys) doAll(item map[string]any, searcher search.Searcher) search.Exporter {
	var contents []string
	for _, _key := range f.keys {
		contents = append(contents, f.GetKeyContent(_key, item))
	}

	return searcher.Do(contents)
}

func (f *Keys) doInOrder(item map[string]any, searcher search.Searcher) search.Exporter {
	var rt search.Exporter

	for _, _key := range f.keys {
		_v := f.GetKeyContent(_key, item)
		rt = searcher.Do([]string{_v})
		if rt.IsOk() {
			break
		}
	}

	return rt
}

func (f *Keys) IsAll() bool {
	return f.way == keysWayAll
}

func (f *Keys) IsInOrder() bool {
	return f.way == keysWayInOrder
}
