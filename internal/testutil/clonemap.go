package testutil

import "testing"

// AssertMapCloned verifies that fn returns an independent map copy each call.
// It mutates the first returned map and checks whether the second is affected.
// Used to verify that DefaultValues() and similar methods properly clone their internal map.
func AssertMapCloned(t *testing.T, name string, fn func() map[string]any) {
	t.Helper()
	dv1 := fn()
	if dv1 == nil {
		return
	}
	dv2 := fn()
	if dv2 == nil {
		return
	}
	dv1["__test_clone__"] = "modified"
	if _, ok := dv2["__test_clone__"]; ok {
		t.Errorf("%s: DefaultValues() returned the same map reference, expected a cloned copy", name)
	}
}
