package testutil

import "testing"

// AssertMapCloned 验证 fn 每次返回独立的 map 副本。
// 修改第一次返回的 map，检查第二次返回的是否受影响。
// 用于验证 DefaultValues() 等方法是否正确克隆了内部 map。
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
		t.Errorf("%s: DefaultValues() 返回了同一个 map 引用，期望返回克隆副本", name)
	}
}
