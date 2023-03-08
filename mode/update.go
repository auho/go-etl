package mode

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/means"
)

// Update
// handle some keys of data for update
type Update struct {
	Mode
	meanses []means.UpdateMeans
}

func NewUpdate(keys []string, meanses ...means.UpdateMeans) *Update {
	m := &Update{}
	m.keys = keys
	m.meanses = meanses

	return m
}

func (u *Update) GetTitle() string {
	is := make([]string, 0)
	for _, i := range u.meanses {
		is = append(is, i.GetTitle())
	}

	return fmt.Sprintf("Update %s{%s}", u.Mode.getTitle(), strings.Join(is, ","))
}

func (u *Update) GetFields() []string {
	return u.keys
}

func (u *Update) Do(item map[string]any) map[string]any {
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

func (u *Update) Close() {
	for k := range u.meanses {
		u.meanses[k].Close()
	}
}
