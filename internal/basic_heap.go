package internal

// NEED TO BE IMPLEMENTED WITH POINTER
// ex: func (s *MyHeap) IsBetterThan(heaper BasicHeaper) {...}
type BasicHeaper interface {
	// if compare like this < heaper -> lowest element popped from heap
	// if compare like this > heaper -> highest element popped first from heap
	IsBetterThan(heaper interface{}) bool

	// Setter and getter for heap index in order to be able to correctly use
	// heap.Fix -> mandatory if heap element value is changed
	SetHeapIndex(index int)
	GetHeapIndex() int
}

type BasicHeap []BasicHeaper

func (h BasicHeap) Len() int           { return len(h) }
func (h BasicHeap) Less(i, j int) bool { return h[i].IsBetterThan(h[j]) }
func (h BasicHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].SetHeapIndex(i)
	h[j].SetHeapIndex(j)
}
func (h *BasicHeap) Push(x interface{}) {
	n := len(*h)
	item := x.(BasicHeaper)
	item.SetHeapIndex(n)
	*h = append(*h, item)
}
func (h *BasicHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	old[n-1] = nil
	x.SetHeapIndex(-1)
	*h = old[0 : n-1]
	return x
}
