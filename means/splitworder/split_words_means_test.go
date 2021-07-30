package splitworder

import (
	"testing"
)

func TestSplitWordsMeans(t *testing.T) {
	content := "1,2,3,4,5,6,7,8,9"

	s := NewSplitWordsMeans(",")
	items := s.Insert([]string{content, content})
	if len(items) != 18 {
		t.Error("error")
	}
}
