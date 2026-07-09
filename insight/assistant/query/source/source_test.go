package source

import (
	"testing"
)

func TestItemValuesToIdentification(t *testing.T) {
	b := &Base{Name: "test", HasNamePrefix: false}
	id := b.itemValuesToIdentification([]string{"a", "b"})
	if id != "a_b" {
		t.Fatalf("expect[a_b] != actual[%s]", id)
	}
}

func TestItemValuesToIdentificationWithPrefix(t *testing.T) {
	b := &Base{Name: "test", HasNamePrefix: true}
	id := b.itemValuesToIdentification([]string{"a", "b"})
	if id != "test_a_b" {
		t.Fatalf("expect[test_a_b] != actual[%s]", id)
	}
}

func TestKeysToIdentification(t *testing.T) {
	b := &Base{}
	id := b.keysToIdentification([]string{"x", "y", "z"})
	if id != "x_y_z" {
		t.Fatalf("expect[x_y_z] != actual[%s]", id)
	}
}
