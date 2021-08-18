package go_etl

import (
	"testing"
)

func TestPython2Go(t *testing.T) {
	p := NewPython2Go("office/const.py")
	p.Conversion()
}
