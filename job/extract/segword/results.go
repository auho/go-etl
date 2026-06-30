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

func (r *result) toTag(format format) map[string]any {
	return map[string]any{
		format.tokenName: r.token,
		format.flagName:  r.flag,
	}
}

type results []result

func (rs results) toAll(format format) []map[string]any {
	var results []map[string]any
	for _, result := range rs {
		results = append(results, result.toTag(format))
	}

	return results
}

func (rs results) toLine(format format) []map[string]any {
	var ss []string

	for _, result := range rs {
		ss = append(ss, result.token)
	}

	return []map[string]any{{format.tokenName: strings.Join(ss, format.sep)}}
}
