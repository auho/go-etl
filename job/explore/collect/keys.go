package collect

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/job/explore/search"
)

var _ Collector = (*Keys)(nil)

const (
	keysWayAll   = iota // 所有的 key
	keysWayFirst        // first
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

// NewKeysFirst
// first matched key, if matched return
func NewKeysFirst(keys []string) *Keys {
	return newKeys(keys, keysWayFirst)
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

func (f *Keys) Do(item map[string]any, searcher search.Searcher) search.Token {
	if f.IsAll() {
		return f.doAll(item, searcher)
	} else if f.IsFirst() {
		return f.doFirst(item, searcher)
	} else {
		panic("way unknown")
	}
}

func (f *Keys) doAll(item map[string]any, searcher search.Searcher) search.Token {
	var contents []string
	for _, _key := range f.keys {
		contents = append(contents, f.GetKeyContent(_key, item))
	}

	return searcher.Do(contents)
}

func (f *Keys) doFirst(item map[string]any, searcher search.Searcher) search.Token {
	var st search.Token

	for _, _key := range f.keys {
		_v := f.GetKeyContent(_key, item)
		st = searcher.Do([]string{_v})
		if st.IsOk() {
			break
		}
	}

	return st
}

func (f *Keys) IsAll() bool {
	return f.way == keysWayAll
}

func (f *Keys) IsFirst() bool {
	return f.way == keysWayFirst
}
