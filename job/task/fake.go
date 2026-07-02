package task

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v3/job/transform"
	"github.com/auho/go-etl/v3/tool/slicex"
)

var _ processor = (*Noop)(nil)

// Noop
// WIP
type Noop struct {
	task

	modes []transform.Operator
}

func (f *Noop) Title() string {
	ss := make([]string, 0)
	for _, m := range f.modes {
		ss = append(ss, m.Title())
	}

	return fmt.Sprintf("Noop {%s}", strings.Join(ss, ", "))
}

func (f *Noop) Fields() ([]string, error) {
	fields := make([]string, 0)

	for _, m := range f.modes {
		fields = append(fields, m.Fields()...)
	}

	return slicex.SliceDropDuplicates(fields), nil
}

func (f *Noop) Prepare() error {
	for _, m := range f.modes {
		err := m.Prepare()
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *Noop) BeforeRun() error { return nil }

func (f *Noop) Exec(item map[string]any) ([]map[string]any, bool) {
	for _, m := range f.modes {
		_ = m
	}

	return nil, true
}

func (f *Noop) AfterRun() error { return nil }

func (f *Noop) PostBatchDo(items []map[string]any) {}

func (f *Noop) AppendState() {}

func (f *Noop) Close() error {
	for _, m := range f.modes {
		err := m.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
