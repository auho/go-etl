package mode

import (
	"testing"
)

func TestTransferMode(t *testing.T) {
	newKeyName := "abc"
	m := NewTransfer(db, ruleTableName, map[string]string{keyName: newKeyName}, map[string]interface{}{"abc": 1})
	if len(m.GetFields()) != 1 {
		t.Error("fields error")
	}

	if len(m.GetKeys()) != 2 {
		t.Error("keys error")
	}

	results := m.Do(item)
	if len(results) != 1 {
		t.Error("Do is error")
	}
}
