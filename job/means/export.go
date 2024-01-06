package means

import (
	"maps"
	strings2 "strings"

	"github.com/auho/go-toolkit/farmtools/convert/types/strings"
)

var _ Exporter = (*Export)(nil)

type Context struct {
	Keys []string
}

type Exporter interface {
	GetKeys() []string
	GetDefaultValues() map[string]any
	Insert([]map[string]any) []map[string]any
}

type Export struct {
	keys          []string
	defaultValues map[string]any

	// []string: keys
	// []map[string]any:
	insert func(Context, []map[string]any) []map[string]any
	update func(Context, map[string]any) map[string]any
}

func (e *Export) GetKeys() []string {
	return e.keys
}

func (e *Export) GetDefaultValues() map[string]any {
	return maps.Clone(e.defaultValues)
}

func (e *Export) Insert(sm []map[string]any) []map[string]any {
	return e.insert(
		Context{Keys: e.keys},
		sm,
	)
}

func NewExport(keys []string, df map[string]any,
	insert func(Context, []map[string]any) []map[string]any,
	update func(Context, map[string]any) map[string]any,
) *Export {
	return &Export{
		keys:          keys,
		defaultValues: df,
		insert:        insert,
		update:        update,
	}
}

func NewExportLine(keys []string, df map[string]any) *Export {
	return NewExport(
		keys,
		df,
		func(ctx Context, sm []map[string]any) []map[string]any {
			if sm == nil {
				return nil
			}

			ss := make(map[string][]string, len(keys))

			for _, _m := range sm {
				for _, key := range keys {
					_s, err := strings.FromAny(_m[key])
					if err != nil {
						panic(err)
					}

					ss[key] = append(ss[key], _s)
				}
			}

			newM := make(map[string]any, len(keys))
			for _, key := range keys {
				newM[key] = strings2.Join(ss[key], "|")
			}

			return []map[string]any{newM}
		},
		nil,
	)
}

func NewExportFlag(flag string, keys []string, df map[string]any) *Export {
	exportLine := NewExportLine(keys, df)
	return NewExport(
		keys,
		df,
		func(ctx Context, sm []map[string]any) []map[string]any {
			nsm := exportLine.Insert(sm)
			if nsm == nil {
				return nil
			}

			nsm[0][flag] = 1

			return nsm
		},
		nil,
	)
}
