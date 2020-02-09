package exercise

import (
	"container/heap"
	"fmt"
	"github.com/madjlzz/hashcode-2020/solver"
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

func ComputeFile(file string) {
	data = parseFile(file)
	//fmt.Printf("%v %v\n", max, data)

	channel := make(chan bool)
	go launchBatch(channel)

	select {
	case <-channel:
		fmt.Println("All possibilities were explored: this should be the best solution")
	case <-time.After(10000 * time.Millisecond):
		fmt.Println("Timedout...")
	}

	best := bestFinished
	if len(solutionsTree) > 0 {
		bestInTree := solutionsTree.Pop()
		if bestInTree != nil && bestInTree.(*Proposition) != nil && bestInTree.(*Proposition).score > best.score {
			best = bestInTree.(*Proposition)
		}
	}

	fmt.Printf("%v\n", best)
}

func launchBatch(channel chan<- bool) {
	start := New([]int{}, 0, 0)
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
	proposition := &Proposition{pizzas, 0, baseScore, 0, parentSwapFrom, 0, true}
	return proposition
}

func newIteration() *Proposition {
	bestProp := heap.Pop(&solutionsTree).(*Proposition)
	if bestProp.score == max {
		return bestProp
	}

	children := bestProp.createChildren()
	if bestProp.canHaveMoreChild {
		heap.Push(&solutionsTree, bestProp)
	} else if bestFinished == nil || bestProp.score > bestFinished.score {
		bestFinished = bestProp
	}

	res := solver.BatchExecution(children, FinalizeProposition, 100000)
	for _, v := range res {
		if v.Res == nil {
			continue
		}
		child := v.Res.(*Proposition)
		if child == nil {
			continue
		}
		if child.score == max {
			return child
		}
	}
	return nil
}

func FinalizeProposition(p interface{}) solver.ExecutionRes {
	if p == nil {
		return solver.ExecutionRes{nil, nil}
	}
	proposition := p.(*Proposition)
	if proposition == nil {
		return solver.ExecutionRes{nil, nil}
	}
	proposition.FullFill()
	proposition.lastElement = proposition.pizzas[len(proposition.pizzas)-1]
	proposition.lastSwapTo = proposition.lastElement
	heap.Push(&solutionsTree, proposition)
	return solver.ExecutionRes{proposition, nil}
}

// This method is used to generate a new proposition from an existing one
// As of now it is very basic meaning bad perfs
// -> given a proposition, a new one will be attempted by switching one of the highest number of parts with the next lowest one
// -> ex  pizzas possible 10 8 6 3 1 / current set [10 6 3] -> 'replace' pizza 6 by 1 -> new set [10 3 1]
func NewChild(p *Proposition) (*Proposition, error) {
	if p.lastElement >= sizeMinus1 {
		p.canHaveMoreChild = false
		return nil, fmt.Errorf("last element reached, no more child possible")
	}

	p.lastSwapTo++
	if p.lastSwapTo < size {
		return New(swapPizzas(p), p.lastSwapFrom, p.score-data[p.pizzas[p.lastSwapFrom]]+data[p.lastSwapTo]), nil
	}

	p.lastSwapFrom++
	if p.lastSwapFrom >= len(p.pizzas) {
		p.canHaveMoreChild = false
		return nil, fmt.Errorf("last element reached, no more child possible")
	}

	p.lastSwapTo = p.lastElement + 1
	return New(swapPizzas(p), p.lastSwapFrom, p.score-data[p.pizzas[p.lastSwapFrom]]+data[p.lastSwapTo]), nil
}

func swapPizzas(p *Proposition) []int {
	newPizzas := make([]int, len(p.pizzas))
	copy(newPizzas[:p.lastSwapFrom], p.pizzas[:p.lastSwapFrom])
	copy(newPizzas[p.lastSwapFrom:], p.pizzas[p.lastSwapFrom+1:])
	newPizzas[len(newPizzas)-1] = p.lastSwapTo
	return newPizzas
}

func (p *Proposition) createChildren() []interface{} {
	res := make([]interface{}, maxChildrenPerLoop)
	var err error
	for i := 0; i < maxChildrenPerLoop; i++ {
		res[i], err = NewChild(p)
		if err != nil {
			break
		}
	}
	return res
}

// Fill with pizzas until it throws up
func (p *Proposition) FullFill() { // -> might be the most costly part of the algo > to be put in separate threads using batchExecution ?
	var start int
	if len(p.pizzas) == 0 {
		start = 0
	} else {
		start = p.pizzas[len(p.pizzas)-1] + 1
	}

	for i := start; i < size; i++ {
		temp := p.score + data[i]
		if temp > max {
			break
		} else {
			p.pizzas = append(p.pizzas, i)
			p.score = temp
		}
	}
}

func parseFile(file string) []int64 {
	content := solver.ReadInput(file)
	size, _ = strconv.Atoi(strings.TrimSpace(content[0][1]))
	res := make([]int64, size)
	sizeMinus1 = size - 1
	for i := 0; i < size; i++ {
		res[i], _ = strconv.ParseInt(strings.TrimSpace(content[1][sizeMinus1-i]), 10, 64)
	}
	max, _ = strconv.ParseInt(content[0][0], 10, 64)
	return res
}
