package task

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/job/transform"
	slices "github.com/auho/go-etl/v2/tool/slicex"
)

var _ processor = (*Fake)(nil)

// Fake
// WIP
type Fake struct {
	task

	modes []transform.Operator
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

func (f *Fake) BeforeRun() error { return nil }

func (f *Fake) Exec(item map[string]any) ([]map[string]any, bool) {
	for _, m := range f.modes {
		_ = m
	}

	return nil, true
}

func (f *Fake) AfterRun() error { return nil }

func (f *Fake) PostBatchDo(items []map[string]any) {}

func (f *Fake) AppendState() {}

func (f *Fake) Close() error {
	for _, m := range f.modes {
		err := m.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
