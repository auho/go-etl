package mode

import (
	"testing"
)

func TestTransferMode(t *testing.T) {
	newKeyName := "abc"
	m := NewTransferMode(ruleTableName, db, map[string]string{keyName: newKeyName})
	for _, key := range m.GetKeys() {
		if key != newKeyName {
			t.Error("keys error")
		}
	}

	for _, field := range m.GetFields() {
		if field != keyName {
			t.Error("keys error")
		}
	}

	results := m.Do(item)
	if len(results) != 1 {
		t.Error("Do is error")
	}
}
