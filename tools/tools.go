package tools

import "testing"

func UnittestAssert(t *testing.T, want, got interface{}) {
	if want != got {
		t.Errorf("Unexpected result: wanted %v got %v", want, got)
	}
}
