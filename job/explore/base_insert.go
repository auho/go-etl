package explore

import (
	"fmt"
	"maps"
	"strings"

	"github.com/auho/go-etl/v2/tool/slices"
)

type baseInsert struct {
	base

	name          string
	is            []*Insert
	insertKeys    []string
	defaultValues map[string]any
}

func (bi *baseInsert) GetTitle() string {
	var ss []string
	for _, _i := range bi.is {
		ss = append(ss, _i.GetTitle())
	}

	return bi.genTitle(bi.name, strings.Join(ss, ","))
}

func (bi *baseInsert) GetFields() []string {
	return bi.keys
}

func (bi *baseInsert) GetKeys() []string {
	var keys []string

	for _, _i := range bi.is {
		keys = append(keys, _i.GetKeys()...)
	}

	keys = slices.SliceDropDuplicates(keys)

	return keys
}

func (bi *baseInsert) DefaultValues() map[string]any {
	return maps.Clone(bi.defaultValues)
}

func (bi *baseInsert) State() []string {
	return []string{fmt.Sprintf("%s: %s", bi.GetTitle(), bi.genCounter())}
}

func (bi *baseInsert) Prepare() error {
	bi.defaultValues = make(map[string]any)

	for _, _i := range bi.is {
		var err error
		err = _i.Prepare()
		if err != nil {
			return err
		}

		bi.keys = append(bi.keys, _i.GetFields()...)
		maps.Copy(bi.defaultValues, _i.DefaultValues())
	}

	return nil
}

func (bi *baseInsert) Close() error {
	var err error
	for _, _i := range bi.is {
		err = _i.Close()
		if err != nil {
			return fmt.Errorf("close err: %w", err)
		}
	}

	return nil
}
