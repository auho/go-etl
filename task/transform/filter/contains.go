package filter

import (
	"fmt"
	strings2 "strings"

	"github.com/auho/go-toolkit/v2/farmtools/convert/types/strings"
)

// NewContainsAll builds a predicate that reports true when the value at key
// contains every string in subs.
func NewContainsAll(key string, subs []string) Predicate {
	return Func(func(m map[string]any) (bool, error) {
		s, err := strings.FromAny(m[key])
		if err != nil {
			return false, fmt.Errorf("FromAny[%s]: %w", key, err)
		}

		for _, sub := range subs {
			if !strings2.Contains(s, sub) {
				return false, nil
			}
		}

		return true, nil
	})
}

// NewContainsAny builds a predicate that reports true when the value at key
// contains any of the strings in subs.
func NewContainsAny(key string, subs []string) Predicate {
	return Func(func(m map[string]any) (bool, error) {
		s, err := strings.FromAny(m[key])
		if err != nil {
			return false, fmt.Errorf("FromAny[%s]: %w", key, err)
		}

		for _, sub := range subs {
			if strings2.Contains(s, sub) {
				return true, nil
			}
		}

		return false, nil
	})
}
