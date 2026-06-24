package transform

import (
	"fmt"
	"maps"

	"github.com/auho/go-etl/v2/job/extract"
)

var _ InsertOperator = (*Insert)(nil)

// Insert
// single means
type Insert struct {
	Mode
	means extract.Inserter
}

func NewInsert(keys []string, means extract.Inserter) *Insert {
	im := &Insert{}
	im.keys = keys
	im.means = means

	return im
}

func (im *Insert) Prepare() error {
	if len(im.keys) <= 0 {
		return fmt.Errorf("Insert Prepare keys not exists error")
	}

	err := im.means.Prepare()
	if err != nil {
		return fmt.Errorf("Insert Prepare error; %w", err)
	}

	return nil
}

func (im *Insert) Title() string {
	return im.GenTitle("Insert", im.means.Title())
}

func (im *Insert) GetFields() []string {
	return im.keys
}

func (im *Insert) Keys() []string {
	return im.means.Keys()
}

func (im *Insert) DefaultValues() map[string]any {
	return maps.Clone(im.means.DefaultValues())
}

func (im *Insert) Do(item map[string]any) []map[string]any {
	im.AddTotal(1)

	if item == nil {
		return nil
	}

	contents := im.GetKeysContent(im.keys, item)
	if len(contents) <= 0 {
		return nil
	}

	rt := im.means.Insert(contents)
	im.AddAmount(int64(len(rt)))

	return rt
}

func (im *Insert) State() []string {
	return []string{fmt.Sprintf("%s: %s", im.Title(), im.GenCounter())}
}

func (im *Insert) Close() error {
	err := im.means.Close()
	if err != nil {
		return fmt.Errorf("Insert close error; %w", err)
	}

	return nil
}
