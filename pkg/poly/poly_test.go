package poly

import (
	"gitlab.doc.ic.ac.uk/js6317/bgw/pkg/field"
	"testing"
)

func TestPoly_Eval(t *testing.T) {
	fld := field.New(101)
	po := New([]int{20, 57, 68}, fld)
	wants := []int{20, 44, 2, 96, 23, 86, 83}
	for j, want := range wants {
		if got := po.Eval(j); got != want {
			t.Errorf("po.Eval = %d, want %d", got, want)
		}
	}
}
