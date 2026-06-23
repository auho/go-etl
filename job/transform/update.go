package transform

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/job/extract"
)

var _ UpdateOperator = (*UpdateMode)(nil)

// UpdateMode
// handle some keys of data for update
type UpdateMode struct {
	Mode
	ms []extract.Updater
}

func NewUpdate(keys []string, ms ...extract.Updater) *UpdateMode {
	um := &UpdateMode{}
	um.keys = keys
	um.ms = ms

	return um
}

func (um *UpdateMode) Prepare() error {
	if len(um.keys) <= 0 {
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

func (um *UpdateMode) Title() string {
	is := make([]string, 0)
	for _, i := range um.ms {
		is = append(is, i.Title())
	}

	return um.GenTitle("UpdateMode", strings.Join(is, ","))
}

func (um *UpdateMode) GetFields() []string {
	return um.keys
}

func (um *UpdateMode) Do(item map[string]any) map[string]any {
	if item == nil {
		return nil
	}

	contents := um.GetKeysContent(um.keys, item)

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
