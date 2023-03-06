package insight

import (
	"testing"
)

func TestPython2Go(t *testing.T) {
	p := NewPython2Go("office/const.py")
	p.clean()
	p.conversionQuote()
	p.conversionComment()
	p.conversionVar()
	p.conversionSlice()
	p.conversionDict()

	p.content()
}
