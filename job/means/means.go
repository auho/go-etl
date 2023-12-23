package means

import (
	"fmt"

	"github.com/auho/go-etl/v2/job/explore/search"
)

var _ InsertMeans = (*Means)(nil)
var _ UpdateMeans = (*Means)(nil)

type Means struct {
	search search.Searcher

	keys          []string
	defaultValues map[string]any
}

func NewMeans(s search.Searcher) *Means {
	m := &Means{search: s}

	return m
}

func (m *Means) Prepare() error {
	err := m.search.Prepare()
	if err != nil {
		return err
	}

	_export := m.search.GenExport()
	m.keys = _export.GetKeys()
	m.defaultValues = _export.GetDefaultValues()

	return nil
}

func (m *Means) GetTitle() string {
	return fmt.Sprintf("means:%s ", m.search.GetTitle())
}

func (m *Means) GetKeys() []string {
	return m.keys
}

func (m *Means) DefaultValues() map[string]any {
	return m.defaultValues
}

func (m *Means) Insert(contents []string) []map[string]any {
	token := m.search.Do(contents)
	rets := token.ToToken()
	if len(rets) <= 0 {
		return nil
	}

	return rets
}

func (m *Means) Update(contents []string) map[string]any {
	token := m.search.Do(contents)
	rets := token.ToToken()
	if len(rets) <= 0 {
		return nil
	}

	return rets[0]
}

func (m *Means) Close() error {
	return m.search.Close()
}
