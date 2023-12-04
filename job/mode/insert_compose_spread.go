package mode

import (
	"fmt"
	"maps"
	slices2 "slices"
	"strings"

	"github.com/auho/go-etl/v2/tool/slices"
)

var _ InsertModer = (*InsertComposeSpreadMode)(nil)

// InsertComposeSpreadMode
// compose spread 取第一个 spread
type InsertComposeSpreadMode struct {
	Mode
	modes []InsertModer

	insertKeys    []string
	defaultValues map[string]any
}

func NewInsertComposeSpread(modes ...InsertModer) *InsertComposeSpreadMode {
	ic := &InsertComposeSpreadMode{}
	ic.modes = modes

	return ic
}

func (ic *InsertComposeSpreadMode) GetTitle() string {
	var ss []string
	for _, m := range ic.modes {
		ss = append(ss, m.GetTitle())
	}

	return ic.GenTitle("InsertComposeSpreadMode", strings.Join(ss, ";"))
}

func (ic *InsertComposeSpreadMode) GetFields() []string {
	for _, m := range ic.modes {
		ic.Keys = append(ic.Keys, m.GetFields()...)
	}

	ic.Keys = slices.SliceDropDuplicates(ic.Keys)

	return slices2.Clone(ic.Keys)
}

func (ic *InsertComposeSpreadMode) GetKeys() []string {
	return ic.insertKeys
}

func (ic *InsertComposeSpreadMode) DefaultValues() map[string]any {
	return maps.Clone(ic.defaultValues)
}

func (ic *InsertComposeSpreadMode) Prepare() error {
	ic.defaultValues = make(map[string]any)

	for _, m := range ic.modes {
		err := m.Prepare()
		if err != nil {
			return err
		}

		ic.insertKeys = append(ic.insertKeys, m.GetKeys()...)

		maps.Copy(ic.defaultValues, m.DefaultValues())
	}

	ic.insertKeys = slices.SliceDropDuplicates(ic.insertKeys)

	return nil
}

func (ic *InsertComposeSpreadMode) Do(item map[string]any) []map[string]any {
	ic.AddTotal(1)

	_has := false
	ret := make(map[string]any)
	for _, m := range ic.modes {
		_mrt := m.Do(item)
		if len(_mrt) <= 0 {
			maps.Copy(ret, m.DefaultValues())
		} else {
			_has = true
			maps.Copy(ret, _mrt[0])
		}
	}

	if _has {
		ic.AddAmount(1)

		return []map[string]any{ret}
	} else {
		return nil
	}
}

func (ic *InsertComposeSpreadMode) State() []string {
	var ss []string
	ss = append(ss, fmt.Sprintf("InsertComposeSpreadMode: %s", ic.GenCounter()))
	for i, m := range ic.modes {
		var mss []string
		for _i, _ms := range m.State() {
			_s := ""
			if _i == 0 {
				_s = fmt.Sprintf("%-5s%s", fmt.Sprintf("%d.", i), _ms)
			} else {
				_s = fmt.Sprintf("%-5s%s", "", _ms)
			}

			mss = append(mss, _s)
		}

		ss = append(ss, mss...)
	}

	return ss
}

func (ic *InsertComposeSpreadMode) Close() error {
	for _, m := range ic.modes {
		err := m.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
