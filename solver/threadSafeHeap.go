package solver

import "container/heap"

/*
Thanks to https://husobee.github.io/heaps/golang/safe/2016/09/01/safe-heaps-golang.html
for the code -> sole modification was to make WatchHeapOps public as as far as I understand, it must be called first by
the app, so that the thread can work properly

Usage:
Main thread start   -> `quit := WatchHeapOps()`
Sub Thread Push     -> `HeapPush(h, &object)`
Sub/Main Thread Pop -> `HeapPop(h)`
Main thread end     -> `quit <- true`
*/

// heapPopChanMsg - the message structure for a pop chan
type heapPopChanMsg struct {
	h      heap.Interface
	result chan interface{}
}

// heapPushChanMsg - the message structure for a push chan
type heapPushChanMsg struct {
	h heap.Interface
	x interface{}
}

var (
	quitChan chan bool
	// heapPushChan - push channel for pushing to a heap
	heapPushChan = make(chan heapPushChanMsg)
	// heapPopChan - pop channel for popping from a heap
	heapPopChan = make(chan heapPopChanMsg)
)

// HeapPush - safely push item to a heap interface
func HeapPush(h heap.Interface, x interface{}) {
	heapPushChan <- heapPushChanMsg{
		h: h,
		x: x,
	}
}

// HeapPop - safely pop item from a heap interface
func HeapPop(h heap.Interface) interface{} {
	var result = make(chan interface{})
	heapPopChan <- heapPopChanMsg{
		h:      h,
		result: result,
	}
	return <-result
}

//stopWatchHeapOps - stop watching for heap operations
func stopWatchHeapOps() {
	quitChan <- true
}

// WatchHeapOps - watch for push/pops to our heap, and serializing the operations with channels
// Must be called at the start to allow pop and push to work in the threads
func WatchHeapOps() chan bool {
	var quit = make(chan bool)
	go func() {
		for {
			select {
			case <-quit:
				return
			case popMsg := <-heapPopChan:
				popMsg.result <- heap.Pop(popMsg.h)
			case pushMsg := <-heapPushChan:
				heap.Push(pushMsg.h, pushMsg.x)
			}
		}
	}()
	return quit
}
