package regexps

import (
	"regexp"

	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
)

var _ search.Searcher = (*Regexp)(nil)

type Regexp struct {
	rule        means.Ruler
	expressions []string
	subMode     func([]*regexp.Regexp, []string) []string
	export      *Export

	regexps []*regexp.Regexp
}

func NewRegexp(rule means.Ruler, exs []string, subMode func([]*regexp.Regexp, []string) []string, export *Export) *Regexp {
	return &Regexp{
		rule:        rule,
		expressions: exs,
		subMode:     subMode,
		export:      export,
	}
}

func (r *Regexp) GetTitle() string {
	return "Regexp"
}

func (r *Regexp) Prepare() error {
	for _, ex := range r.expressions {
		r.regexps = append(r.regexps, regexp.MustCompile(ex))
	}

	return nil
}

func (r *Regexp) GenExport() search.Exporter {
	return r.export
}

func (r *Regexp) Do(contents []string) search.Token {
	rets := r.subMode(r.regexps, contents)

	return r.export.ToToken(rets, r.rule)
}

func (r *Regexp) Close() error { return nil }

func (r *Regexp) ToMeans() *means.Means {
	return means.NewMeans(r)
}

func NewRegexpAll(rule means.Ruler, exs []string, export *Export) *Regexp {
	return NewRegexp(rule, exs, func(regexps []*regexp.Regexp, contents []string) []string {
		var rets []string
		for _, content := range contents {
			for _, re := range regexps {
				ret := re.FindStringSubmatch(content)
				if ret != nil {
					rets = append(rets, ret[1])
				}
			}
		}

		return nil
	}, export)
}

func NewRegexpFirst(rule means.Ruler, exs []string, export *Export) *Regexp {
	return NewRegexp(rule, exs, func(regexps []*regexp.Regexp, contents []string) []string {
		var rets []string
		for _, content := range contents {
			for _, re := range regexps {
				ret := re.FindStringSubmatch(content)
				if ret != nil {
					rets = append(rets, ret[1])
					goto LOOP
				}
			}
		}
	LOOP:

		return rets
	}, export)
}
