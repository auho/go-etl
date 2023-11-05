package mode

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/job/means"
)

var _ UpdateModer = (*UpdateMode)(nil)

// UpdateMode
// handle some keys of data for update
type UpdateMode struct {
	Mode
	meanses []means.UpdateMeans
}

func NewUpdateMode(keys []string, meanses ...means.UpdateMeans) *UpdateMode {
	um := &UpdateMode{}
	um.keys = keys
	um.meanses = meanses

	return um
}

func (um *UpdateMode) Prepare() error {
	for _, m := range um.meanses {
		err := m.Prepare()
		if err != nil {
			return fmt.Errorf("update prepare error; %w", err)
		}
	}

	return nil
}

func (um *UpdateMode) GetTitle() string {
	is := make([]string, 0)
	for _, i := range um.meanses {
		is = append(is, i.GetTitle())
	}

	return fmt.Sprintf("UpdateMode %s{%s}", um.Mode.getTitle(), strings.Join(is, ","))
}

func (um *UpdateMode) GetFields() []string {
	return um.keys
}

func (um *UpdateMode) Do(item map[string]any) map[string]any {
	if item == nil {
		return nil
	}

	contents := um.GetKeysContent(um.keys, item)

	m := make(map[string]any)
	for _, uMeans := range um.meanses {
		_m := uMeans.Update(contents)
		for _k, _v := range _m {
			m[_k] = _v
		}
	}

	return m
}

func (um *UpdateMode) Close() error {
	for k := range um.meanses {
		err := um.meanses[k].Close()
		if err != nil {
			return fmt.Errorf("UpdateMode close error; %w", err)
		}
	}

	return nil
}
