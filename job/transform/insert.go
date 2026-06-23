package transform

import (
	"fmt"
	"maps"

	"github.com/auho/go-etl/v2/job/extract"
)

var _ InsertOperator = (*InsertMode)(nil)

// InsertMode
// single means
type InsertMode struct {
	Mode
	means extract.Inserter
}

func NewInsert(keys []string, means extract.Inserter) *InsertMode {
	im := &InsertMode{}
	im.Keys = keys
	im.means = means

	return im
}

func (im *InsertMode) Prepare() error {
	if len(im.Keys) <= 0 {
		return fmt.Errorf("InsertMode Prepare keys not exists error")
	}

	err := im.means.Prepare()
	if err != nil {
		return fmt.Errorf("InsertMode Prepare error; %w", err)
	}

	return nil
}

func (im *InsertMode) Title() string {
	return im.GenTitle("InsertMode", im.means.Title())
}

func (im *InsertMode) GetFields() []string {
	return im.Keys
}

func (im *InsertMode) GetKeys() []string {
	return im.means.GetKeys()
}

func (im *InsertMode) DefaultValues() map[string]any {
	return maps.Clone(im.means.DefaultValues())
}

func (im *InsertMode) Do(item map[string]any) []map[string]any {
	im.AddTotal(1)

	if item == nil {
		return nil
	}

	contents := im.GetKeysContent(im.Keys, item)
	if len(contents) <= 0 {
		return nil
	}

	rt := im.means.Insert(contents)
	im.AddAmount(int64(len(rt)))

	return rt
}

func (im *InsertMode) State() []string {
	return []string{fmt.Sprintf("%s: %s", im.Title(), im.GenCounter())}
}

func (im *InsertMode) Close() error {
	err := im.means.Close()
	if err != nil {
		return fmt.Errorf("InsertMode close error; %w", err)
	}

	return nil
}
