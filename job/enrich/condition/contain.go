package condition

import (
	strings2 "strings"

	"github.com/auho/go-toolkit/farmtools/convert/types/strings"
)

// NewContainAll
// key contain all subs
func NewContainAll(key string, subs []string) Operation {
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

// NewContainAny
// key contain any one subs
func NewContainAny(key string, subs []string) Operation {
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
