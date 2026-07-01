package segword

import (
	"strings"
)

const nameToken = "token"
const nameFlag = "flag"

var defaultFormat = format{
	tokenName: nameToken,
	flagName:  nameFlag,
	sep:       " ",
}

type format struct {
	tokenName string
	flagName  string
	sep       string
}

func (f *format) check() {
	if f.tokenName == "" {
		f.tokenName = nameToken
	}

	if f.flagName == "" {
		f.flagName = nameFlag
	}

	if f.sep == "" {
		f.sep = ""
	}
}

type result struct {
	token string
	flag  string
}

func (r *result) toTag(f format) map[string]any {
	return map[string]any{
		f.tokenName: r.token,
		f.flagName:  r.flag,
	}
}

type results []result

func (rs results) toAll(f format) []map[string]any {
	var rets []map[string]any
	for _, ret := range rs {
		rets = append(rets, ret.toTag(f))
	}

	return rets
}

func (rs results) toLine(f format) []map[string]any {
	var ss []string

	for _, ret := range rs {
		ss = append(ss, ret.token)
	}

	return []map[string]any{{f.tokenName: strings.Join(ss, f.sep)}}
}
