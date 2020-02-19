package iterativeExecution

import (
	"container/heap"
	"fmt"
	"github.com/MadJlzz/hashcode-2020/algo/batchExecution"
	"github.com/MadJlzz/hashcode-2020/internal"
	"reflect"
	"time"
)

const (
	ItExGlobalTimeout = 60000
	ItExLocalTimeout  = 10000
)

type IterativeExecutioner interface {
	// Create N children from current instances (can be nil, but will then be ignored. Operations can be not thread safe
	CreateChildren() []interface{}

	// Will stop trying to create children from it if false
	CanHaveMoreChild() bool

	// Check if instance is better than input
	IsBetterThan(interface{}) bool

	// Compute the *score* of the element and any other heavy calculations. Only THREAD SAFE operations are allowed here
	Compute()

	// Define if instance reached the goal
	IsMax() bool
}

var solutionsTree internal.BasicHeap
var bestFinished IterativeExecutioner

func IterativeExecution(startItem IterativeExecutioner) IterativeExecutioner {
	bestFinished = nil
	solutionsTree = nil

	start := time.Now()

	quit := internal.WatchHeapOps()
	channel := make(chan bool)
	go launchBatch(startItem, channel)

	select {
	case <-channel:
		fmt.Println("Jackpot: Best solution found !")
	case <-time.After(ItExGlobalTimeout * time.Millisecond):
		fmt.Println("Timedout...")
	}

	best := bestFinished
	if len(solutionsTree) > 0 {
		bestInTree := heap.Pop(&solutionsTree)
		if bestInTree != nil && !reflect.ValueOf(bestInTree).IsNil() && bestInTree.(IterativeExecutioner).IsBetterThan(best) {
			best = bestInTree.(IterativeExecutioner)
		}
	}

	fmt.Printf("Took %s\n", time.Since(start))
	quit <- true

	return best
}

func launchBatch(startItem IterativeExecutioner, quit chan<- bool) {
	startItem.Compute()
	heap.Push(&solutionsTree, startItem)

	for i := 0; len(solutionsTree) > 0; i++ {
		newIteration(quit)
	}
	quit <- true
}

// Get best proposition from heap, create N children from it by making small modification to its solution
// Fullfill its children in a batch execution. Stop all if one of the child has the max score
func newIteration(quit chan<- bool) {
	bestProp := heap.Pop(&solutionsTree).(IterativeExecutioner)
	if bestProp.IsMax() {
		bestFinished = bestProp
		quit <- true
	}

	children := bestProp.CreateChildren()
	if bestProp.CanHaveMoreChild() {
		heap.Push(&solutionsTree, bestProp)
	} else if bestProp.IsBetterThan(bestFinished) {
		bestFinished = bestProp
	}

	batchExecution.BatchExecutionBasic(children, iterativeExecutionCompute, ItExLocalTimeout)
}

func iterativeExecutionCompute(o interface{}) interface{} {
	var emptyRes IterativeExecutioner
	if o == nil || reflect.ValueOf(o).IsNil() {
		return emptyRes
	}
	object := o.(IterativeExecutioner)

	object.Compute()
	channel := internal.HeapPushWithFeedback(&solutionsTree, object)

	// below is a hack used to force waiting for the end of the push to release the thread
	<-channel
	return object
}
