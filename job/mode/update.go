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
	ms []means.UpdateMeans
}

func NewUpdate(keys []string, ms ...means.UpdateMeans) *UpdateMode {
	um := &UpdateMode{}
	um.Keys = keys
	um.ms = ms

	return um
}

func (um *UpdateMode) Prepare() error {
	if len(um.Keys) <= 0 {
		return fmt.Errorf("update prepare keys is not exists error")
	}

	if len(um.ms) <= 0 {
		return fmt.Errorf("update prepare ms error")
	}

	for _, m := range um.ms {
		err := m.Prepare()
		if err != nil {
			return fmt.Errorf("update prepare error; %w", err)
		}
	}

	return nil
}

func (um *UpdateMode) GetTitle() string {
	is := make([]string, 0)
	for _, i := range um.ms {
		is = append(is, i.GetTitle())
	}

	return um.GenTitle("UpdateMode", strings.Join(is, ","))
}

func (um *UpdateMode) GetFields() []string {
	return um.Keys
}

func (um *UpdateMode) Do(item map[string]any) map[string]any {
	if item == nil {
		return nil
	}

	contents := um.GetKeysContent(um.Keys, item)

	if len(contents) <= 0 {
		return nil
	}

	m := make(map[string]any)
	for _, uMeans := range um.ms {
		_m := uMeans.Update(contents)
		for _k, _v := range _m {
			m[_k] = _v
		}
	}

	return m
}

func (um *UpdateMode) Close() error {
	for k := range um.ms {
		err := um.ms[k].Close()
		if err != nil {
			return fmt.Errorf("UpdateMode close error; %w", err)
		}
	}

	return nil
}
