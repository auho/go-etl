package tag

import (
	"fmt"
	"strings"

	"github.com/auho/go-etl/v2/job/explore/search"
	"github.com/auho/go-etl/v2/job/means"
)

var _ search.Searcher = (*SearchKeyword)(nil)

type SearchKeyword struct {
	tagSearch
	genExport GenExportKeyword
	fn        func(*SearchKeyword, []string) search.Exporter
}

func NewSearchKeyword(rule means.Ruler, gek GenExportKeyword, fn func(*SearchKeyword, []string) search.Exporter) *SearchKeyword {
	return &SearchKeyword{
		tagSearch: tagSearch{rule: rule},
		genExport: gek,
		fn:        fn,
	}
}

func (s *SearchKeyword) GetTitle() string {
	return fmt.Sprintf("SearchKeyword{%s}", strings.Join(s.GenExport().GetKeys(), ","))
}

func (s *SearchKeyword) GenExport() search.Exporter {
	return s.genExport(nil, s.rule)
}

func (s *SearchKeyword) Do(contents []string) search.Exporter {
	return s.fn(s, contents)
}

func newSearchKeyword(rule means.Ruler, gek GenExportKeyword, genDoFn func(*SearchKeyword) func([]string) Results) *SearchKeyword {
	return NewSearchKeyword(rule, gek, func(s *SearchKeyword, contents []string) search.Exporter {
		rets := genDoFn(s)(contents)
		if rets == nil {
			return nil
		}

		return s.genExport(rets, rule)
	})
}

func NewSearchKey(rule means.Ruler, gek GenExportKeyword) search.Searcher {
	return newSearchKeyword(rule, gek, func(s *SearchKeyword) func([]string) Results {
		return s.matcher.MatchKey
	})
}

func NewSearchFirstKey(rule means.Ruler, gek GenExportKeyword) search.Searcher {
	return newSearchKeyword(rule, gek, func(s *SearchKeyword) func([]string) Results {
		return s.matcher.MatchFirstKey
	})
}

func NewSearchFirstText(rule means.Ruler, gek GenExportKeyword) search.Searcher {
	return newSearchKeyword(rule, gek, func(s *SearchKeyword) func([]string) Results {
		return s.matcher.MatchFirstText
	})
}

func NewSearchMostKey(rule means.Ruler, gek GenExportKeyword) search.Searcher {
	return newSearchKeyword(rule, gek, func(s *SearchKeyword) func([]string) Results {
		return s.matcher.MatchMostKey
	})
}

func NewSearchMostText(rule means.Ruler, gek GenExportKeyword) search.Searcher {
	return newSearchKeyword(rule, gek, func(s *SearchKeyword) func([]string) Results {
		return s.matcher.MatchMostText
	})
}
