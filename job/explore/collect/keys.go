package collect

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/job/explore/token"
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

func (f *Keys) Pick(item map[string]any, fn func([]string) token.Tokenizer) token.Tokenizer {
	if f.IsAll() {
		return f.doAll(item, fn)
	} else if f.IsInOrder() {
		return f.doInOrder(item, fn)
	} else {
		panic("way unknown")
	}
}

func (f *Keys) doAll(item map[string]any, fn func([]string) token.Tokenizer) token.Tokenizer {
	var rt token.Tokenizer

	var contents []string
	for _, _key := range f.keys {
		contents = append(contents, f.GetKeyContent(_key, item))
	}

	rt = fn(contents)

	return rt
}

func (f *Keys) doInOrder(item map[string]any, fn func([]string) token.Tokenizer) token.Tokenizer {
	var rt token.Tokenizer
	for _, _key := range f.keys {
		_v := f.GetKeyContent(_key, item)
		rt = fn([]string{_v})
		if rt.GetOk() {
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
