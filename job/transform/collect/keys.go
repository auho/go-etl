package collect

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v3/job/extract"
)

var _ Collector = (*Keys)(nil)

const (
	modeAll = iota // collect all keys
	modeAny        // collect just any one
)

// Keys
// collect from keys
type Keys struct {
	ContentReader

	keys []string
	mode int
}

// NewKeys
// collect all keys
func NewKeys(keys []string) *Keys {
	return newKeys(keys, modeAll)
}

// NewKeysAny
// collect any one, if matched return
func NewKeysAny(keys []string) *Keys {
	return newKeys(keys, modeAny)
}

func newKeys(keys []string, mode int) *Keys {
	return &Keys{
		keys: keys,
		mode: mode,
	}
}

func (k *Keys) Title() string {
	return fmt.Sprintf("keys{%s}", strings.Join(k.keys, ","))
}

func (k *Keys) SourceKeys() []string {
	return k.keys
}

func (k *Keys) Extract(item map[string]any, e extract.Extractor) (extract.Result, error) {
	if k.IsAll() {
		return k.doAll(item, e)
	} else if k.IsAny() {
		return k.doAny(item, e)
	} else {
		panic("mode unknown")
	}
}

func (k *Keys) doAll(item map[string]any, e extract.Extractor) (extract.Result, error) {
	var contents []string
	for _, _key := range k.keys {
		content, err := k.Content(_key, item)
		if err != nil {
			return extract.Result{}, fmt.Errorf("Keys.doAll: %w", err)
		}

		contents = append(contents, content)
	}

	return e.Search(contents), nil
}

func (k *Keys) doAny(item map[string]any, e extract.Extractor) (extract.Result, error) {
	var st extract.Result

	for _, _key := range k.keys {
		_v, err := k.Content(_key, item)
		if err != nil {
			return extract.Result{}, fmt.Errorf("Content: %w", err)
		}

		st = e.Search([]string{_v})
		if st.IsOK() {
			break
		}
	}

	return st, nil
}

func (k *Keys) IsAll() bool {
	return k.mode == modeAll
}

func (k *Keys) IsAny() bool {
	return k.mode == modeAny
}
