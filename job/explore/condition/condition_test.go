package condition

import (
	"testing"
)

var _ Conditioner = (*condInt)(nil)
var _ Conditioner = (*condString)(nil)

type condInt struct {
	int int
}

func (c *condInt) OK(m map[string]any) bool {
	return m["int"] == c.int
}

type condString struct {
	string string
}

func (c *condString) OK(m map[string]any) bool {
	return m["string"] == c.string
}

func TestNewAND(t *testing.T) {
	_item := map[string]any{
		"int":    1,
		"string": "1",
	}

	var a AND

	a = NewAND(&condInt{int: 1}, &condString{string: "1"})
	if !a.OK(_item) {
		t.Fatal()
	}

	a = NewAND(&condInt{int: 1}, &condString{string: "2"})
	if a.OK(_item) {
		t.Fatal()
	}

	a = NewAND(&condInt{int: 2}, &condString{string: "1"})
	if a.OK(_item) {
		t.Fatal()
	}

	a = NewAND(&condInt{int: 2}, &condString{string: "2"})
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

	o = NewOR(&condInt{int: 1}, &condString{string: "1"})
	if !o.OK(_item) {
		t.Fatal()
	}

	o = NewOR(&condInt{int: 1}, &condString{string: "2"})
	if !o.OK(_item) {
		t.Fatal()
	}

	o = NewOR(&condInt{int: 2}, &condString{string: "1"})
	if !o.OK(_item) {
		t.Fatal()
	}

	o = NewOR(&condInt{int: 2}, &condString{string: "2"})
	if o.OK(_item) {
		t.Fatal()
	}
}
