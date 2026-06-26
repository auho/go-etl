package transform

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v3/job/extract"
)

var _ UpdateOperator = (*Update)(nil)

// Update
// handle some keys of data for update
type Update struct {
	base
	updaters []extract.Updater
}

func NewUpdate(keys []string, updaters ...extract.Updater) *Update {
	um := &Update{}
	um.keys = keys
	um.updaters = updaters

	return um
}

func (um *Update) Prepare() error {
	if len(um.keys) <= 0 {
		return fmt.Errorf("keys do not exist")
	}

	if len(um.updaters) <= 0 {
		return fmt.Errorf("inserters do not exist")
	}

	for _, m := range um.updaters {
		err := m.Prepare()
		if err != nil {
			return fmt.Errorf("prepare: %w", err)
		}
	}

	return nil
}

func (um *Update) Title() string {
	is := make([]string, 0)
	for _, i := range um.updaters {
		is = append(is, i.Title())
	}

	return um.GenTitle("Update", strings.Join(is, ","))
}

func (um *Update) GetFields() []string {
	return um.keys
}

func (um *Update) Do(item map[string]any) map[string]any {
	if item == nil {
		return nil
	}

	contents := um.GetKeysContent(um.keys, item)

	if len(contents) <= 0 {
		return nil
	}

	m := make(map[string]any)
	for _, uMeans := range um.updaters {
		_m := uMeans.Update(contents)
		for _k, _v := range _m {
			m[_k] = _v
		}
	}

	return m
}

func (um *Update) Close() error {
	for k := range um.updaters {
		err := um.updaters[k].Close()
		if err != nil {
			return fmt.Errorf("close: %w", err)
		}
	}

	return nil
}
