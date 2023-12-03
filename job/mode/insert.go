package mode

import (
	"fmt"
	"maps"

	"github.com/auho/go-etl/v2/job/means"
)

var _ InsertModer = (*InsertMode)(nil)

// InsertMode
// single means
type InsertMode struct {
	Mode
	means means.InsertMeans
}

func NewInsert(keys []string, means means.InsertMeans) *InsertMode {
	im := &InsertMode{}
	im.Keys = keys
	im.means = means

	return im
}

func (im *InsertMode) Prepare() error {
	if len(im.Keys) <= 0 {
		return fmt.Errorf("InsertMode Prepare kyes not exists error")
	}

	err := im.means.Prepare()
	if err != nil {
		return fmt.Errorf("InsertMode Prepare error; %w", err)
	}

	return nil
}

func (im *InsertMode) GetTitle() string {
	return im.GenTitle("InsertMode", im.means.GetTitle())
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
	if item == nil {
		return nil
	}

	contents := im.GetKeysContent(im.Keys, item)
	if len(contents) <= 0 {
		return nil
	}

	return im.means.Insert(contents)
}

func (im *InsertMode) Close() error {
	err := im.means.Close()
	if err != nil {
		return fmt.Errorf("InsertMode close error; %w", err)
	}

	return nil
}
