package transform

import (
	"fmt"
	"maps"

	"github.com/auho/go-etl/v3/job/extract"
)

var _ InsertOperator = (*Insert)(nil)

// Insert
// single inserter
type Insert struct {
	base
	inserter extract.Inserter
}

func NewInsert(keys []string, inserter extract.Inserter) *Insert {
	im := &Insert{}
	im.keys = keys
	im.inserter = inserter

	return im
}

func (im *Insert) Prepare() error {
	if len(im.keys) <= 0 {
		return fmt.Errorf("keys do not exist")
	}

	err := im.inserter.Prepare()
	if err != nil {
		return fmt.Errorf("inserter.Prepare: %w", err)
	}

	return nil
}

func (im *Insert) Title() string {
	return im.GenTitle("Insert", im.inserter.Title())
}

func (im *Insert) GetFields() []string {
	return im.keys
}

func (im *Insert) Keys() []string {
	return im.inserter.Keys()
}

func (im *Insert) DefaultValues() map[string]any {
	return maps.Clone(im.inserter.DefaultValues())
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

	rt := im.inserter.Insert(contents)
	im.AddAmount(int64(len(rt)))

	return rt
}

func (im *Insert) State() []string {
	return []string{fmt.Sprintf("%s: %s", im.Title(), im.GenCounter())}
}

func (im *Insert) Close() error {
	err := im.inserter.Close()
	if err != nil {
		return fmt.Errorf("inserter.Close: %w", err)
	}

	return nil
}
