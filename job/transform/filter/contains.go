package filter

import (
	"fmt"
	strings2 "strings"

	"github.com/auho/go-toolkit/v2/farmtools/convert/types/strings"
)

// NewContainsAll
// key contains all subs
func NewContainsAll(key string, subs []string) Predicate {
	return func(m map[string]any) (bool, error) {
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
	}
}

// NewContainsAny
// key contains any of subs
func NewContainsAny(key string, subs []string) Predicate {
	return func(m map[string]any) (bool, error) {
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
	}
}
