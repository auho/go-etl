package match

import (
	"fmt"

	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
)

var _ means.InsertMeans = (*Means)(nil)
var _ means.UpdateMeans = (*Means)(nil)

// Means
// match means
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
	return fmt.Sprintf("Match:%s ", m.search.GetTitle())
}

func (m *Means) GetKeys() []string {
	return m.keys
}

func (m *Means) DefaultValues() map[string]any {
	return m.defaultValues
}

func (m *Means) Insert(contents []string) []map[string]any {
	return m.search.Do(contents).ToTokenize()
}

func (m *Means) Update(contents []string) map[string]any {
	results := m.search.Do(contents).ToTokenize()
	if results == nil {
		return nil
	}

	return results[0]
}

func (m *Means) Close() error {
	return m.search.Close()
}

// NewWholeLabels
// merge all labels together
// label1|label2|label3
// keyword1|keyword2|keyword3|
func NewWholeLabels(rule means.Ruler) *Means {
	return NewMeans(NewSearchWholeLabels(rule))
}

// NewLabel
// label tags
func NewLabel(rule means.Ruler) *Means {
	return NewMeans(NewSearchLabels(rule, NewExportLabelAll))
}

// NewKey
// keyword
func NewKey(rule means.Ruler) *Means {
	return NewMeans(NewSearchKey(rule, NewExportKeywordAll))
}

// NewFirstKey
// the first matched key
func NewFirstKey(rule means.Ruler) *Means {
	return NewMeans(NewSearchFirstKey(rule, NewExportKeywordAll))
}

// NewFirstText
// the leftmost text
func NewFirstText(rule means.Ruler) *Means {
	return NewMeans(NewSearchFirstText(rule, NewExportKeywordAll))
}

// NewMostKey
// most key
func NewMostKey(rule means.Ruler) *Means {
	return NewMeans(NewSearchMostKey(rule, NewExportKeywordAll))
}

// NewMostText
// most text
func NewMostText(rule means.Ruler) *Means {
	return NewMeans(NewSearchMostText(rule, NewExportKeywordAll))
}
