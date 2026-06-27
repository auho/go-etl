package transform

import (
	"fmt"
	"maps"
	slices2 "slices"
	"strings"

	slices "github.com/auho/go-etl/v3/tool/slicex"
)

var _ InsertOperator = (*InsertComposeSpread)(nil)

// InsertComposeSpread
// compose spread 取第一个 spread
type InsertComposeSpread struct {
	base
	operators []InsertOperator

	insertKeys    []string
	defaultValues map[string]any
}

func NewInsertComposeSpread(operators ...InsertOperator) *InsertComposeSpread {
	ic := &InsertComposeSpread{}
	ic.operators = operators

	return ic
}

func (ic *InsertComposeSpread) Title() string {
	var ss []string
	for _, m := range ic.operators {
		ss = append(ss, m.Title())
	}

	return ic.GenTitle("InsertComposeSpread", strings.Join(ss, ";"))
}

func (ic *InsertComposeSpread) GetFields() []string {
	for _, m := range ic.operators {
		ic.keys = append(ic.keys, m.GetFields()...)
	}

	ic.keys = slices.SliceDropDuplicates(ic.keys)

	return slices2.Clone(ic.keys)
}

func (ic *InsertComposeSpread) Keys() []string {
	return ic.insertKeys
}

func (ic *InsertComposeSpread) DefaultValues() map[string]any {
	return maps.Clone(ic.defaultValues)
}

func (ic *InsertComposeSpread) Prepare() error {
	ic.defaultValues = make(map[string]any)

	for _, m := range ic.operators {
		err := m.Prepare()
		if err != nil {
			return err
		}

		ic.insertKeys = append(ic.insertKeys, m.Keys()...)

		maps.Copy(ic.defaultValues, m.DefaultValues())
	}

	ic.insertKeys = slices.SliceDropDuplicates(ic.insertKeys)

	return nil
}

func (ic *InsertComposeSpread) Apply(item map[string]any) ([]map[string]any, error) {
	ic.AddTotal(1)

	_has := false
	ret := make(map[string]any)
	for _, m := range ic.operators {
		_mrt, err := m.Apply(item)
		if err != nil {
			return nil, fmt.Errorf("apply: %w", err)
		}
		if len(_mrt) <= 0 {
			maps.Copy(ret, m.DefaultValues())
		} else {
			_has = true
			maps.Copy(ret, _mrt[0])
		}
	}

	if _has {
		ic.AddAmount(1)

		return []map[string]any{ret}, nil
	} else {
		return nil, nil
	}
}

func (ic *InsertComposeSpread) State() []string {
	var ss []string
	ss = append(ss, fmt.Sprintf("InsertComposeSpread: %s", ic.GenCounter()))
	for i, m := range ic.operators {
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

func (ic *InsertComposeSpread) Close() error {
	for _, m := range ic.operators {
		err := m.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
