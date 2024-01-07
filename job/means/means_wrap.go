package means

var _ InsertMeans = (*Wrap)(nil)
var _ UpdateMeans = (*Wrap)(nil)

type Wrap struct {
	means  *Means
	export Exporter

	keys          []string
	defaultValues map[string]any
	hasExport     bool
}

// NewWrap
// Deprecated
func NewWrap(m *Means, export *Export) *Wrap {
	return &Wrap{
		means:  m,
		export: export,
	}
}

func (w *Wrap) GetTitle() string {
	return w.means.GetTitle()
}

func (w *Wrap) GetKeys() []string {
	return w.keys
}

func (w *Wrap) DefaultValues() map[string]any {
	return w.defaultValues
}

func (w *Wrap) Prepare() error {
	err := w.means.Prepare()
	if err != nil {
		return err
	}

	w.keys = w.means.GetKeys()
	w.defaultValues = w.means.DefaultValues()
	if w.export != nil {
		w.hasExport = true
		w.keys = w.export.GetKeys()
		w.defaultValues = w.export.GetDefaultValues()
	}

	return nil
}

func (w *Wrap) Insert(contents []string) []map[string]any {
	rets := w.means.Insert(contents)
	if rets == nil {
		return nil
	}

	if w.hasExport {
		rets = w.export.Insert(rets)
	}

	return rets
}

func (w *Wrap) Update(contents []string) map[string]any {
	rets := w.means.Insert(contents)
	if rets == nil {
		return nil
	}

	return rets[0]
}

func (w *Wrap) Close() error {
	return w.means.Close()
}
