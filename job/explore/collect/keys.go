package collect

import (
	"strings"

	"github.com/auho/go-etl/v2/job/explore/token"
)

var _ Collector = (*Keys)(nil)

// Keys
// collect from keys
type Keys struct {
	Collect

	keys []string
}

func NewKeys(keys []string) *Keys {
	f := &Keys{}
	f.keys = keys

	return f
}

func (f *Keys) GetTitle() string {
	return strings.Join(f.keys, ",")
}

func (f *Keys) GetKeys() []string {
	return f.keys
}

func (f *Keys) Search(item map[string]any, fn func(string) token.Tokenizer) token.Tokenizer {
	var rt token.Tokenizer
	for _, _key := range f.keys {
		_v := f.GetKeyContent(_key, item)
		rt = fn(_v)
		if rt.GetOk() {
			break
		}
	}

	return rt
}
