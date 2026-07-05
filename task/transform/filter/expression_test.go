package filter

import (
	"errors"
	"testing"
)

var opInt = func(i int) Predicate {
	return Func(func(m map[string]any) (bool, error) {
		return m["int"] == i, nil
	})
}

var opString = func(s string) Predicate {
	return Func(func(m map[string]any) (bool, error) {
		return m["string"] == s, nil
	})
}

// errorPredicate always fails on Prepare.
type errorPredicate struct{}

func (errorPredicate) Prepare() error                         { return errors.New("prepare failed") }
func (errorPredicate) Match(map[string]any) (bool, error)     { return false, nil }

func TestNewAND(t *testing.T) {
	_item := map[string]any{
		"int":    1,
		"string": "1",
	}

	var a And

	a = NewAnd(opInt(1), opString("1"))
	ok, err := a.Match(_item)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal()
	}

	a = NewAnd(opInt(1), opString("2"))
	ok, err = a.Match(_item)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal()
	}

	a = NewAnd(opInt(2), opString("1"))
	ok, err = a.Match(_item)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal()
	}

	a = NewAnd(opInt(2), opString("2"))
	ok, err = a.Match(_item)
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

	var o Or

	o = NewOr(opInt(1), opString("1"))
	ok, err := o.Match(_item)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal()
	}

	o = NewOr(opInt(1), opString("2"))
	ok, err = o.Match(_item)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal()
	}

	o = NewOr(opInt(2), opString("1"))
	ok, err = o.Match(_item)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal()
	}

	o = NewOr(opInt(2), opString("2"))
	ok, err = o.Match(_item)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal()
	}
}

func TestExpression_Match(t *testing.T) {
	// And.Match: true only when all operands match
	andPred := NewAnd(opInt(1), opString("1"))
	ok, err := andPred.Match(map[string]any{"int": 1, "string": "1"})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("And: expected true for all match")
	}
	ok, err = andPred.Match(map[string]any{"int": 2, "string": "1"})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("And: expected false when first operand fails")
	}
	ok, err = andPred.Match(map[string]any{"int": 1, "string": "2"})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("And: expected false when second operand fails")
	}
	ok, err = andPred.Match(map[string]any{"int": 2, "string": "2"})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("And: expected false when all operands fail")
	}

	// Or.Match: true when any operand matches
	orPred := NewOr(opInt(1), opString("2"))
	ok, err = orPred.Match(map[string]any{"int": 1, "string": "1"})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("Or: expected true when first operand matches")
	}
	ok, err = orPred.Match(map[string]any{"int": 2, "string": "2"})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("Or: expected true when second operand matches")
	}
	ok, err = orPred.Match(map[string]any{"int": 1, "string": "2"})
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("Or: expected true when both operands match")
	}
	ok, err = orPred.Match(map[string]any{"int": 2, "string": "3"})
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("Or: expected false when no operand matches")
	}
}

func TestExpression_Prepare(t *testing.T) {
	// And.Prepare / Or.Prepare propagate errors from operands.
	errPred := errorPredicate{}
	if err := NewAnd(opInt(1), errPred).Prepare(); err == nil {
		t.Fatal("And.Prepare: expected error from failing operand")
	}
	if err := NewOr(errPred, opInt(1)).Prepare(); err == nil {
		t.Fatal("Or.Prepare: expected error from failing operand")
	}

	// Successful operands return nil.
	if err := NewAnd(opInt(1), opString("1")).Prepare(); err != nil {
		t.Fatalf("And.Prepare: expected nil, got %v", err)
	}
	if err := NewOr(opInt(1), opString("1")).Prepare(); err != nil {
		t.Fatalf("Or.Prepare: expected nil, got %v", err)
	}
}
