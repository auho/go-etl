package keys

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v3/task/transform/collector/source"
	typesStrings "github.com/auho/go-toolkit/v2/farmtools/convert/types/strings"
)

var _ source.Source = (*Keys)(nil)

// Keys is a Source that reads content from a map by string keys.
type Keys struct {
	keys []string
}

func New(ks []string) *Keys {
	return &Keys{keys: ks}
}

func (k *Keys) Title() string {
	return fmt.Sprintf("keys{%s}", strings.Join(k.keys, ","))
}

func (k *Keys) Keys() []string {
	return k.keys
}

func (k *Keys) Prepare() error {
	if len(k.keys) == 0 {
		return fmt.Errorf("keys is empty")
	}
	return nil
}

// Contents reads values for all keys from item.
// A missing key defaults to empty string (no error).
// A type conversion failure (FromAny) is still an error.
func (k *Keys) Contents(item map[string]any) ([]string, map[string]string, error) {
	keysValue := make(map[string]string, len(k.keys))
	for _, key := range k.keys {
		raw, ok := item[key]
		if !ok {
			keysValue[key] = ""
			continue
		}
		v, err := typesStrings.FromAny(raw)
		if err != nil {
			return nil, nil, fmt.Errorf("FromAny[%s]%T: %w", key, raw, err)
		}
		keysValue[key] = v
	}
	return k.keys, keysValue, nil
}
