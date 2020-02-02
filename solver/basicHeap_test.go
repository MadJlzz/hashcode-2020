package solver

import (
	"container/heap"
	"testing"
)

type intHeap struct {
	a, index int
}

var useLess = true

func (i *intHeap) Compare(heaper BasicHeaper) bool {
	if useLess {
		return i.a < heaper.(*intHeap).a
	} else {
		return i.a > heaper.(*intHeap).a
	}
}

func (i *intHeap) SetHeapIndex(index int) {
	i.index = index
}

func (i *intHeap) GetHeapIndex() int {
	return i.index
}

func TestBasicHeap_int(t *testing.T) {
	useLess = true
	a1 := &intHeap{a: 10}
	a2 := &intHeap{a: 1}
	a3 := &intHeap{a: 20}
	basicHeap := &BasicHeap{a1}
	heap.Init(basicHeap)

	heap.Push(basicHeap, a2)
	heap.Push(basicHeap, a3)

	Assert(t, 3, basicHeap.Len())
	Assert(t, 1, heap.Pop(basicHeap).(*intHeap).a)
	Assert(t, 10, heap.Pop(basicHeap).(*intHeap).a)
	Assert(t, 20, heap.Pop(basicHeap).(*intHeap).a)
}

func TestBasicHeap_intReversed(t *testing.T) {
	useLess = false
	a1 := &intHeap{a: 10}
	a2 := &intHeap{a: 1}
	a3 := &intHeap{a: 20}
	basicHeap := &BasicHeap{a1}
	heap.Init(basicHeap)

	heap.Push(basicHeap, a2)
	heap.Push(basicHeap, a3)

	Assert(t, 3, basicHeap.Len())
	Assert(t, 20, heap.Pop(basicHeap).(*intHeap).a)
	Assert(t, 10, heap.Pop(basicHeap).(*intHeap).a)
	Assert(t, 1, heap.Pop(basicHeap).(*intHeap).a)
}

func TestBasicHeap_intChanged(t *testing.T) {
	useLess = true
	a1 := &intHeap{a: 10}
	a2 := &intHeap{a: 1}
	a3 := &intHeap{a: 20}
	basicHeap := &BasicHeap{a1}
	heap.Init(basicHeap)

	heap.Push(basicHeap, a2)
	heap.Push(basicHeap, a3)

	a2.a = 100
	heap.Fix(basicHeap, a2.index)

	Assert(t, 3, basicHeap.Len())
	Assert(t, 10, heap.Pop(basicHeap).(*intHeap).a)
	Assert(t, 20, heap.Pop(basicHeap).(*intHeap).a)
	Assert(t, 100, heap.Pop(basicHeap).(*intHeap).a)
}

func Assert(t *testing.T, want, got interface{}) {
	if want != got {
		t.Errorf("Unexpected result: wanted %v got %v", want, got)
	}
}
