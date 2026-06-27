package collect

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v3/job/extract"
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

func (k *Keys) Title() string {
	return fmt.Sprintf("keys{%s}", strings.Join(k.keys, ","))
}

func (k *Keys) Keys() []string {
	return k.keys
}

func (k *Keys) Search(item map[string]any, searcher extract.Extractor) (extract.Result, error) {
	if k.IsAll() {
		return k.doAll(item, searcher)
	} else if k.IsAny() {
		return k.doAny(item, searcher)
	} else {
		panic("way unknown")
	}
}

func (k *Keys) doAll(item map[string]any, searcher extract.Extractor) (extract.Result, error) {
	var contents []string
	for _, _key := range k.keys {
		content, err := k.GetKeyContent(_key, item)
		if err != nil {
			return extract.Result{}, fmt.Errorf("Keys.doAll: %w", err)
		}

		contents = append(contents, content)
	}

	return searcher.Search(contents), nil
}

func (k *Keys) doAny(item map[string]any, searcher extract.Extractor) (extract.Result, error) {
	var st extract.Result

	for _, _key := range k.keys {
		_v, err := k.GetKeyContent(_key, item)
		if err != nil {
			return extract.Result{}, fmt.Errorf("GetKeyContent: %w", err)
		}

		st = searcher.Search([]string{_v})
		if st.IsOK() {
			break
		}
	}

	return st, nil
}

func (k *Keys) IsAll() bool {
	return k.way == keysWayAll
}

func (k *Keys) IsAny() bool {
	return k.way == keysWayAny
}
