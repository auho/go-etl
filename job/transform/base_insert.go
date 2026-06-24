package transform

import (
	"fmt"
	"maps"
	"strings"

	"github.com/auho/go-etl/v2/job/extract"
)

// insertHorizontal
// 多个 means horizontal
type insertHorizontal struct {
	Mode
	ms []extract.Inserter

	insertKeys    []string
	defaultValues map[string]any
}

func newInsertHorizontal(keys []string, ms ...extract.Inserter) insertHorizontal {
	ih := insertHorizontal{}
	ih.keys = keys
	ih.ms = ms

	return ih
}

func (ih *insertHorizontal) Prepare() error {
	if len(ih.keys) <= 0 {
		return fmt.Errorf("insertHorizontal Prepare keys not exists error")
	}

	for _, m := range ih.ms {
		err := m.Prepare()
		if err != nil {
			return fmt.Errorf("prepare error; %w", err)
		}
	}

	ih.defaultValues = make(map[string]any)

	for _, m := range ih.ms {
		ih.insertKeys = append(ih.insertKeys, m.Keys()...)

		maps.Copy(ih.defaultValues, m.DefaultValues())
	}

	return nil
}

func (ih *insertHorizontal) Title() string {
	var ss []string
	for _, m := range ih.ms {
		ss = append(ss, m.Title())
	}

	return ih.GenTitle("insertHorizontal", strings.Join(ss, ","))
}

func (ih *insertHorizontal) GetFields() []string {
	return ih.keys
}

func (ih *insertHorizontal) Keys() []string {
	return ih.insertKeys
}

func (ih *insertHorizontal) DefaultValues() map[string]any {
	return maps.Clone(ih.defaultValues)
}

func (ih *insertHorizontal) State() []string {
	return []string{fmt.Sprintf("%s: %s", ih.Title(), ih.GenCounter())}
}

func (ih *insertHorizontal) Close() error {
	for _, m := range ih.ms {
		err := m.Close()
		if err != nil {
			return fmt.Errorf("close error; %w", err)
		}
	}

	return nil
}
