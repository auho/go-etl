package mode

import (
	"fmt"
	"maps"
	"strings"

	"github.com/auho/go-etl/v2/job/means"
)

// insertHorizontalMode
// 多个 means horizontal
type insertHorizontalMode struct {
	Mode
	ms []means.InsertMeans

	insertKeys    []string
	defaultValues map[string]any
}

func newInsertHorizontal(keys []string, ms ...means.InsertMeans) insertHorizontalMode {
	ih := insertHorizontalMode{}
	ih.Keys = keys
	ih.ms = ms

	return ih
}

func (ih *insertHorizontalMode) Prepare() error {
	if len(ih.Keys) <= 0 {
		return fmt.Errorf("insertHorizontalMode Prepare kyes not exists error")
	}

	for _, m := range ih.ms {
		err := m.Prepare()
		if err != nil {
			return fmt.Errorf("prepare error; %w", err)
		}
	}

	ih.defaultValues = make(map[string]any)

	for _, m := range ih.ms {
		ih.insertKeys = append(ih.insertKeys, m.GetKeys()...)

		maps.Copy(ih.defaultValues, m.DefaultValues())
	}

	return nil
}

func (ih *insertHorizontalMode) GetTitle() string {
	var ss []string
	for _, m := range ih.ms {
		ss = append(ss, m.GetTitle())
	}

	return ih.GenTitle("insertHorizontalMode", strings.Join(ss, ","))
}

func (ih *insertHorizontalMode) GetFields() []string {
	return ih.Keys
}

func (ih *insertHorizontalMode) GetKeys() []string {
	return ih.insertKeys
}

func (ih *insertHorizontalMode) DefaultValues() map[string]any {
	return maps.Clone(ih.defaultValues)
}

func (ih *insertHorizontalMode) Close() error {
	for _, m := range ih.ms {
		err := m.Close()
		if err != nil {
			return fmt.Errorf("close error; %w", err)
		}
	}

	return nil
}
