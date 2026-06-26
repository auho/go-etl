package filter

import (
	strings2 "strings"

	"github.com/auho/go-toolkit/v2/farmtools/convert/types/strings"
)

// NewContainsAll
// key contains all subs
func NewContainsAll(key string, subs []string) Predicate {
	return func(m map[string]any) bool {
		s, err := strings.FromAny(m[key])
		if err != nil {
			panic(err)
		}

		for _, sub := range subs {
			if !strings2.Contains(s, sub) {
				return false
			}
		}

		return true
	}
}

// NewContainsAny
// key contains any of subs
func NewContainsAny(key string, subs []string) Predicate {
	return func(m map[string]any) bool {
		s, err := strings.FromAny(m[key])
		if err != nil {
			panic(err)
		}

		for _, sub := range subs {
			if strings2.Contains(s, sub) {
				return true
			}
		}

		return false
	}
}
