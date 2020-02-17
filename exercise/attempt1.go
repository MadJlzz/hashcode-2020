package exercise

import (
	"container/heap"
	"fmt"
	"github.com/madjlzz/hashcode-2020/solver"
	"sort"
	"strconv"
	"strings"
	"time"
)

var max int64
var data []int64
var size int
var sizeMinus1 int

const maxChildrenPerLoop = 100

var solutionsTree solver.BasicHeap
var bestFinished *Proposition

func SolveExercise(fileContent map[int][]string) (res [][]string) {
	bestFinished = nil
	solutionsTree = nil

	start := time.Now()

	data = parseFile(fileContent)

	quit := solver.WatchHeapOps()
	channel := make(chan bool)
	go launchBatch(channel)

	select {
	case <-channel:
		fmt.Println("Jackpot: Best solution found !")
	case <-time.After(100000000 * time.Millisecond):
		fmt.Println("Timedout...")
	}

	best := bestFinished
	if len(solutionsTree) > 0 {
		bestInTree := heap.Pop(&solutionsTree)
		if bestInTree != nil && bestInTree.(*Proposition) != nil && (best == nil || bestInTree.(*Proposition).score > best.score) {
			best = bestInTree.(*Proposition)
		}
	}

	pizzaRes := reorder(best)
	fmt.Printf("score=%v length=%d res=%v\n", best.score, len(best.pizzas), pizzaRes)

	fmt.Printf("Took %s\n", time.Since(start))
	quit <- true

	res = append(res, []string{fmt.Sprintf("%d", len(best.pizzas))})
	res = append(res, make([]string, len(best.pizzas)))
	for i := 0; i < len(best.pizzas); i++ {
		res[1][i] = fmt.Sprintf("%d", pizzaRes)
	}
	return res
}

func reorder(p *Proposition) []int {
	var s = len(data) - 1
	res := make([]int, len(p.pizzas))
	for i := 0; i < len(p.pizzas); i++ {
		res[i] = s - p.pizzas[i]
	}
	sort.Ints(res)
	return res
}

func launchBatch(channel chan<- bool) {
	start := &Proposition{[]int{}, 0, 0, 0, 0, 1, true}
	FinalizeProposition(start)

	var res *Proposition
	for i := 0; res == nil && len(solutionsTree) > 0; i++ {
		//fmt.Printf("Best score at turn %d: %v  -- %v\n", i, solutionsTree[0], bestFinished)
		res = newIteration()
	}
	channel <- true
}

type Res struct {
	slices []int
	res    int
}

type Proposition struct {
	pizzas      []int
	lastElement int

	score     int64
	heapIndex int

	// Used to know which was the last child generated from it. If last swap matches the last possible one,
	// no more child can be made and it is extracted from the heap, compared to the previous bestScore and kept or not depending on the comparison
	// Child is created by removing pizza i and adding instead new smaller pizza j -> pizzas possible 10 8 6 3 1 / current set [10 6 3] -> 'replace' pizza 6 by 1 -> new set [10 3 1] to be FullFill
	lastSwapFrom     int // id of the element in the array Proposition.pizzas
	lastSwapTo       int // id of the element in the data array
	canHaveMoreChild bool
}

func (p *Proposition) Compare(heaper solver.BasicHeaper) bool {
	return p.score > (heaper.(*Proposition)).score
}
func (p *Proposition) SetHeapIndex(index int) { p.heapIndex = index }
func (p *Proposition) GetHeapIndex() int      { return p.heapIndex }

func New(pizzas []int, parentSwapFrom int, baseScore int64) *Proposition {
	proposition := &Proposition{pizzas, 0, baseScore, 0, parentSwapFrom, pizzas[parentSwapFrom] + 1, true}
	return proposition
}

// Get best proposition from heap, create N children from it by making small modification to its solution
// Fullfill its children in a batch execution. Stop all if one of the child has the max score
func newIteration() *Proposition {
	bestProp := heap.Pop(&solutionsTree).(*Proposition)
	if bestProp.score == max {
		bestFinished = bestProp
		return bestProp
	}

	children := bestProp.createChildren()
	if bestProp.canHaveMoreChild {
		heap.Push(&solutionsTree, bestProp)
	} else if bestFinished == nil || bestProp.score > bestFinished.score {
		bestFinished = bestProp
	}

	solver.BatchExecutionBasic(children, FinalizeProposition, 100000000)
	return nil
}

// Most costly part of the algo and thus put in thread
// Fill remaining pizzas until max value is reached, and push to heap
func FinalizeProposition(p interface{}) interface{} {
	var emptyRes *Proposition
	proposition := p.(*Proposition)
	if proposition == nil {
		return emptyRes
	}
	proposition.FullFill()
	proposition.lastElement = proposition.pizzas[len(proposition.pizzas)-1]
	channel := solver.HeapPushWithFeedback(&solutionsTree, proposition)

	// below is a hack used to force waiting for the end of the push to release the thread
	select {
	case <-channel:
	}
	return proposition
}

// This method is used to generate a new proposition from an existing one
// As of now it is very basic meaning bad perfs
// -> given a proposition, a new one will be attempted by switching one of the highest number of parts with the next lowest one
// -> ex  pizzas possible 10 8 6 3 1 / current set [10 6 3] -> 'replace' pizza 6 by 1 -> new set [10 3 1]
func NewChild(i interface{}) interface{} {
	p := i.(*Proposition)
	var emptyRes *Proposition
	if !p.canHaveMoreChild {
		return emptyRes
	}

	if p.lastSwapFrom >= len(p.pizzas) {
		p.canHaveMoreChild = false
		return emptyRes
	}

	if p.lastSwapTo >= sizeMinus1 {
		p.lastSwapFrom++
		if p.lastSwapFrom >= len(p.pizzas) {
			p.canHaveMoreChild = false
			return emptyRes
		}

		p.lastSwapTo = p.pizzas[p.lastSwapFrom] + 1
		return NewChild(p)
	}

	newScore := p.score - data[p.pizzas[p.lastSwapFrom]]
	for p.lastSwapTo <= sizeMinus1 {
		if Contains(p.pizzas, p.lastSwapTo) {
			p.lastSwapTo++
			continue
		}

		temp := newScore + data[p.lastSwapTo]
		if temp <= max {
			res := New(swapPizzas(p), p.lastSwapFrom, temp)
			p.lastSwapTo++
			return res
		}
		p.lastSwapTo++
	}
	return NewChild(p)
}

// to get new base score for child, take current, subtract the value of the removed pizza and add the value of the new one
func (p *Proposition) getNewBaseScore() int64 {
	return p.score - data[p.pizzas[p.lastSwapFrom]] + data[p.pizzas[p.lastSwapTo]]
}

// Remove item at index lastSwapFrom and add item at the end with value lastSwapTo
func swapPizzas(p *Proposition) []int {
	newPizzas := make([]int, len(p.pizzas))
	copy(newPizzas[:p.lastSwapFrom], p.pizzas[:p.lastSwapFrom])
	copy(newPizzas[p.lastSwapFrom:], p.pizzas[p.lastSwapFrom+1:])
	newPizzas[len(newPizzas)-1] = p.lastSwapTo
	return newPizzas
}

// Create n children from current proposition
func (p *Proposition) createChildren() []interface{} {
	res := make([]interface{}, maxChildrenPerLoop)
	for i := 0; i < maxChildrenPerLoop; i++ {
		res[i] = NewChild(p)
	}
	return res
}

// Fill with pizzas until it throws up
func (p *Proposition) FullFill() { // -> might be the most costly part of the algo > to be put in separate threads using batchExecution ?
	var start int
	if len(p.pizzas) == 0 {
		start = 0
	} else {
		start = p.pizzas[p.lastSwapFrom] + 1
	}

	for i := start; i < size; i++ {
		if Contains(p.pizzas, i) {
			continue
		}
		temp := p.score + data[i]
		if temp > max {
			continue
		} else {
			p.pizzas = append(p.pizzas, i)
			p.score = temp
		}
	}
}

func parseFile(fileContent map[int][]string) []int64 {
	size, _ = strconv.Atoi(strings.TrimSpace(fileContent[0][1]))
	res := make([]int64, size)
	sizeMinus1 = size - 1
	for i := 0; i < size; i++ {
		res[i], _ = strconv.ParseInt(strings.TrimSpace(fileContent[1][sizeMinus1-i]), 10, 64)
	}
	max, _ = strconv.ParseInt(fileContent[0][0], 10, 64)

	return res
}

func Contains(a []int, x int) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
