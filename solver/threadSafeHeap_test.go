package solver

import (
	"testing"
)

type threadSafeHeap struct {
	a, index int
}

func (i *threadSafeHeap) IsBetterThan(heaper interface{}) bool {
	return i.a < heaper.(*threadSafeHeap).a
}
func (i *threadSafeHeap) SetHeapIndex(index int) { i.index = index }
func (i *threadSafeHeap) GetHeapIndex() int      { return i.index }

var loopNb = 10000

func TestHeapThreadSafe(t *testing.T) {
	quit := WatchHeapOps()

	channel := make(chan bool)
	basicHeap := &BasicHeap{}

	// Launch N successive push to heap in different thread (would crash if not done using threadSafeHeap)
	for i := loopNb - 1; i >= 0; i-- {
		go pushToHeap(channel, i, basicHeap)
	}

	// Wait for all push to complete
	for i := 0; i < loopNb; i++ {
		select {
		case <-channel:
		}
	}

	// Assert that result is ordered as expected
	for i := 0; i < loopNb; i++ {
		if item := HeapPop(basicHeap).(*threadSafeHeap).a; item != i {
			t.Errorf("Was expecting %d got %d", i, item)
		}
	}
	quit <- true
}

func pushToHeap(channel chan<- bool, i int, h *BasicHeap) {
	HeapPush(h, &threadSafeHeap{a: i})
	channel <- true
}
