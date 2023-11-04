package mode

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/means"
)

var _ UpdateModer = (*UpdateMode)(nil)

// UpdateMode
// handle some keys of data for update
type UpdateMode struct {
	Mode
	meanses []means.UpdateMeans
}

func NewUpdateMode(keys []string, meanses ...means.UpdateMeans) *UpdateMode {
	m := &UpdateMode{}
	m.keys = keys
	m.meanses = meanses

	return m
}

func (u *UpdateMode) Prepare() error {
	for _, m := range u.meanses {
		err := m.Prepare()
		if err != nil {
			return fmt.Errorf("update prepare error; %w", err)
		}
	}

	return nil
}

func (u *UpdateMode) GetTitle() string {
	is := make([]string, 0)
	for _, i := range u.meanses {
		is = append(is, i.GetTitle())
	}

	return fmt.Sprintf("UpdateMode %s{%s}", u.Mode.getTitle(), strings.Join(is, ","))
}

func (u *UpdateMode) GetFields() []string {
	return u.keys
}

func (u *UpdateMode) Do(item map[string]any) map[string]any {
	if item == nil {
		return nil
	}

	contents := u.GetKeysContent(u.keys, item)

	m := make(map[string]any)
	for _, uMeans := range u.meanses {
		_m := uMeans.Update(contents)
		for _k, _v := range _m {
			m[_k] = _v
		}
	}

	return m
}

func (u *UpdateMode) Close() error {
	for k := range u.meanses {
		err := u.meanses[k].Close()
		if err != nil {
			return fmt.Errorf("UpdateMode close error; %w", err)
		}
	}

	return nil
}
