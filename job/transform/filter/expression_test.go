package filter

import (
	"testing"
)

var opInt = func(i int) Predicate {
	return func(m map[string]any) (bool, error) {
		return m["int"] == i, nil
	}
}

var opString = func(s string) Predicate {
	return func(m map[string]any) (bool, error) {
		return m["string"] == s, nil
	}
}

func TestNewAND(t *testing.T) {
	_item := map[string]any{
		"int":    1,
		"string": "1",
	}

	var a AND

	a = NewAND(opInt(1), opString("1"))
	ok, err := a.OK(_item)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal()
	}

	a = NewAND(opInt(1), opString("2"))
	ok, err = a.OK(_item)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal()
	}

	a = NewAND(opInt(2), opString("1"))
	ok, err = a.OK(_item)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal()
	}

	a = NewAND(opInt(2), opString("2"))
	ok, err = a.OK(_item)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
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
	ok, err := o.OK(_item)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal()
	}

	o = NewOR(opInt(1), opString("2"))
	ok, err = o.OK(_item)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal()
	}

	o = NewOR(opInt(2), opString("1"))
	ok, err = o.OK(_item)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal()
	}

	o = NewOR(opInt(2), opString("2"))
	ok, err = o.OK(_item)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal()
	}
}
