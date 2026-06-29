package keys

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v3/job/transform/collector"
	typesStrings "github.com/auho/go-toolkit/v2/farmtools/convert/types/strings"
)

var _ collector.Source = (*Keys)(nil)

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

func (k *Keys) Contents(ks []string, item map[string]any) ([]string, error) {
	contents := make([]string, 0, len(ks))
	for _, key := range ks {
		v, err := typesStrings.FromAny(item[key])
		if err != nil {
			return nil, fmt.Errorf("FromAny[%s]%T: %w", key, item[key], err)
		}
		contents = append(contents, v)
	}
	return contents, nil
}