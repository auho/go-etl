package segword

import (
	"strings"
)

const NameToken = "token"
const NameFlag = "flag"

var DefaultFormat = Format{
	TokenName: NameToken,
	FlagName:  NameFlag,
	Sep:       " ",
}

type Format struct {
	TokenName string
	FlagName  string
	Sep       string
}

func (f *Format) check() {
	if f.TokenName == "" {
		f.TokenName = NameToken
	}

	if f.FlagName == "" {
		f.FlagName = NameFlag
	}

	if f.Sep == "" {
		f.Sep = ""
	}
}

type Result struct {
	Token string
	Flag  string
}

func (r *Result) ToTag(format Format) map[string]any {
	return map[string]any{
		format.TokenName: r.Token,
		format.FlagName:  r.Flag,
	}
}

type Results []Result

func (rs Results) ToAll(format Format) []map[string]any {
	var results []map[string]any
	for _, result := range rs {
		results = append(results, result.ToTag(format))
	}

	return results
}

func (rs Results) ToLine(format Format) []map[string]any {
	var ss []string

	for _, result := range rs {
		ss = append(ss, result.Token)
	}

	return []map[string]any{{format.TokenName: strings.Join(ss, format.Sep)}}
}
