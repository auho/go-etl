package collect

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/job/explore/search"
)

var _ Collector = (*Keys)(nil)

const (
	keysWayAll = iota // collect all keys
	keysWayAny        // collect just any one
)

// Keys
// collect from keys
type Keys struct {
	Collect

	keys []string
	way  int
}

// NewKeys
// collect all keys
func NewKeys(keys []string) *Keys {
	return newKeys(keys, keysWayAll)
}

// NewKeysAny
// collect any one, if matched return
func NewKeysAny(keys []string) *Keys {
	return newKeys(keys, keysWayAny)
}

func newKeys(keys []string, way int) *Keys {
	return &Keys{
		keys: keys,
		way:  way,
	}
}

func (k *Keys) GetTitle() string {
	return fmt.Sprintf("keys{%s}", strings.Join(k.keys, ","))
}

func (k *Keys) GetKeys() []string {
	return k.keys
}

func (k *Keys) Do(item map[string]any, searcher search.Searcher) search.Token {
	if k.IsAll() {
		return k.doAll(item, searcher)
	} else if k.IsAny() {
		return k.doAny(item, searcher)
	} else {
		panic("way unknown")
	}
}

func (k *Keys) doAll(item map[string]any, searcher search.Searcher) search.Token {
	var contents []string
	for _, _key := range k.keys {
		contents = append(contents, k.GetKeyContent(_key, item))
	}

	return searcher.Do(contents)
}

func (k *Keys) doAny(item map[string]any, searcher search.Searcher) search.Token {
	var st search.Token

	for _, _key := range k.keys {
		_v := k.GetKeyContent(_key, item)
		st = searcher.Do([]string{_v})
		if st.IsOk() {
			break
		}
	}

	return st
}

func (k *Keys) IsAll() bool {
	return k.way == keysWayAll
}

func (k *Keys) IsAny() bool {
	return k.way == keysWayAny
}
