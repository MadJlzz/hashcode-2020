package internal

import "testing"

var max float64
var dataFloat []float64

func TestSliceChildren(t *testing.T) {
	SliceChildrenCanAdd = func(a *ArrayScored, i int) (bool, float64) {
		newScore := a.Score + dataFloat[i]
		return newScore < max, newScore
	}

	max = 25
	dataFloat = []float64{1, 2, 4, 6, 8, 10, 12, 14, 16, 18}
	SliceChildrenSize = len(dataFloat)

	start := SliceChildren{&ArrayScored{}, &ArrayScored{}, 0, true}
	a1 := start.GetChild()
	a2 := start.GetChild()
	a3 := start.GetChild()

	start.Complete()
	a1.Complete()
	a2.Complete()
	a3.Complete()

	println(a1, a2, a3)
}
