package action

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/job/mode"
	"github.com/auho/go-etl/v2/tool/slices"
)

var _ Actor = (*Fake)(nil)

// Fake
// WIP
type Fake struct {
	Action

	modes []mode.Moder
}

func (f *Fake) Title() string {
	ss := make([]string, 0)
	for _, m := range f.modes {
		ss = append(ss, m.GetTitle())
	}

	return fmt.Sprintf("Fake {%s}", strings.Join(ss, ", "))
}

func (f *Fake) GetFields() []string {
	fields := make([]string, 0)

	for _, m := range f.modes {
		fields = append(fields, m.GetFields()...)
	}

	return slices.SliceDropDuplicates(fields)
}

func (f *Fake) Prepare() error {
	for _, m := range f.modes {
		err := m.Prepare()
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *Fake) PreDo() error { return nil }

func (f *Fake) Do(item map[string]any) ([]map[string]any, bool) {
	for _, m := range f.modes {
		_ = m
	}

	return nil, true
}

func (f *Fake) PostDo() error { return nil }

func (f *Fake) PostBatchDo(items []map[string]any) {}

func (f *Fake) Blink() {}

func (f *Fake) Close() error {
	for _, m := range f.modes {
		err := m.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
