package solver

import (
	"math"
	"testing"
)

var target float64

var nbChild = 500
var maxChild = 1000

type IterativeExecutionTest struct {
	generation int
	lastChild  int
	nb         float64
	score      float64
	heapIndex  int
}

func (o *IterativeExecutionTest) SetHeapIndex(index int) {
	o.heapIndex = index
}

func (o *IterativeExecutionTest) GetHeapIndex() int {
	return o.heapIndex
}

func (o *IterativeExecutionTest) CreateChildren() []interface{} {
	res := make([]interface{}, nbChild)
	for i := 0; i < nbChild; i++ {
		if !o.CanHaveMoreChild() {
			break
		}
		res[i] = &IterativeExecutionTest{o.generation + 1, 0, o.nb + (o.nb * float64(o.lastChild) / float64(maxChild)), 0, 0}
		o.lastChild++
	}
	return res
}

func (i *IterativeExecutionTest) CanHaveMoreChild() bool {
	return i.lastChild != maxChild && i.nb < target
}

func (i *IterativeExecutionTest) IsBetterThan(o interface{}) bool {
	return o == nil || o.(*IterativeExecutionTest) == nil || i.score > o.(*IterativeExecutionTest).score
}

func (i *IterativeExecutionTest) Compute() {
	if i.nb > target {
		i.score = -1
	} else if i.nb == target {
		i.score = math.MaxFloat64
	} else {
		i.score = 1 / (target - i.nb)
	}
}

func (i *IterativeExecutionTest) IsMax() bool {
	return i.nb == target
}

func TestIterativeExecution(t *testing.T) {
	target = 1.23456789
	start := IterativeExecutionTest{0, 0, 1, 0, 0}
	res := IterativeExecution(&start)
	if res.(*IterativeExecutionTest).nb != target {
		t.Errorf("error")
	}
}
