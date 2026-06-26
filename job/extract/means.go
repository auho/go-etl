package extract

import (
	"fmt"
)

var _ Inserter = (*Means)(nil)
var _ Updater = (*Means)(nil)

type Means struct {
	search Extractor
	export Exporter

	keys          []string
	defaultValues map[string]any
	hasExport     bool
}

// NewMeans
// Deprecated: change to using explore.Searcher
func NewMeans(s Extractor) *Means {
	m := &Means{search: s}

	return m
}

func (m *Means) Prepare() error {
	err := m.search.Prepare()
	if err != nil {
		return err
	}

	_export := m.search.NewExport()
	m.keys = _export.Keys()
	m.defaultValues = _export.DefaultValues()
	if m.export != nil {
		m.hasExport = true
		m.keys = m.export.Keys()
		m.defaultValues = m.export.DefaultValues()
	}

	return nil
}

func (m *Means) Title() string {
	return fmt.Sprintf("means:%s ", m.search.Title())
}

func (m *Means) Keys() []string {
	return m.keys
}

func (m *Means) DefaultValues() map[string]any {
	return m.defaultValues
}

func (m *Means) Insert(contents []string) []map[string]any {
	token := m.search.Search(contents)
	rets := token.Rows()
	if len(rets) <= 0 {
		return nil
	}

	if m.hasExport {
		rets = m.export.Insert(rets)
	}

	return rets
}

func (m *Means) Update(contents []string) map[string]any {
	token := m.search.Search(contents)
	rets := token.Rows()
	if len(rets) <= 0 {
		return nil
	}

	return rets[0]
}

func (m *Means) Close() error {
	return m.search.Close()
}

func (m *Means) WithExport(e *Export) *Means {
	m.export = e

	return m
}
