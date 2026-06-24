package transform

import (
	"fmt"
	"maps"
	"strings"

	"github.com/auho/go-etl/v3/job/extract"
)

// insertHorizontal
// 多个 inserter horizontal
type insertHorizontal struct {
	base
	inserters []extract.Inserter

	insertKeys    []string
	defaultValues map[string]any
}

func newInsertHorizontal(keys []string, inserters ...extract.Inserter) insertHorizontal {
	ih := insertHorizontal{}
	ih.keys = keys
	ih.inserters = inserters

	return ih
}

func (ih *insertHorizontal) Prepare() error {
	if len(ih.keys) <= 0 {
		return fmt.Errorf("insertHorizontal Prepare keys not exists error")
	}

	for _, m := range ih.inserters {
		err := m.Prepare()
		if err != nil {
			return fmt.Errorf("prepare error; %w", err)
		}
	}

	ih.defaultValues = make(map[string]any)

	for _, m := range ih.inserters {
		ih.insertKeys = append(ih.insertKeys, m.Keys()...)

		maps.Copy(ih.defaultValues, m.DefaultValues())
	}

	return nil
}

func (ih *insertHorizontal) Title() string {
	var ss []string
	for _, m := range ih.inserters {
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
	for _, m := range ih.inserters {
		err := m.Close()
		if err != nil {
			return fmt.Errorf("close error; %w", err)
		}
	}

	return nil
}
