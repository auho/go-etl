package expression

import (
	"testing"
)

var opInt = func(i int) Operation {
	return func(m map[string]any) bool {
		return m["int"] == i
	}
}

var opString = func(s string) Operation {
	return func(m map[string]any) bool {
		return m["string"] == s
	}
}

func TestNewAND(t *testing.T) {
	_item := map[string]any{
		"int":    1,
		"string": "1",
	}

	var a AND

	a = NewAND(opInt(1), opString("1"))
	if !a.OK(_item) {
		t.Fatal()
	}

	a = NewAND(opInt(1), opString("2"))
	if a.OK(_item) {
		t.Fatal()
	}

	a = NewAND(opInt(2), opString("1"))
	if a.OK(_item) {
		t.Fatal()
	}

	a = NewAND(opInt(2), opString("2"))
	if a.OK(_item) {
		t.Fatal()
	}
}

func TestNewOR(t *testing.T) {
	_item := map[string]any{
		"int":    1,
		"string": "1",
	}

	var o OR

	o = NewOR(opInt(1), opString("1"))
	if !o.OK(_item) {
		t.Fatal()
	}

	o = NewOR(opInt(1), opString("2"))
	if !o.OK(_item) {
		t.Fatal()
	}

	o = NewOR(opInt(2), opString("1"))
	if !o.OK(_item) {
		t.Fatal()
	}

	o = NewOR(opInt(2), opString("2"))
	if o.OK(_item) {
		t.Fatal()
	}
}
