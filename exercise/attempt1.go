package exercise

import (
	"fmt"
	"github.com/MadJlzz/hashcode-2020/algo/iterativeExecution"
	"sort"
	"strconv"
	"strings"
)

var max int64
var data []int64
var size int
var sizeMinus1 int

const maxChildrenPerLoop = 100

func SolveExercise(fileContent map[int][]string) (res [][]string) {
	data = nil
	data = parseFile(fileContent)

	start := &Proposition{[]int{}, 0, 0, 0, 0, 1, true}

	best := iterativeExecution.IterativeExecution(start).(*Proposition)

	pizzaRes := reorder(best)
	fmt.Printf("score=%v length=%d res=%v\n", best.score, len(best.pizzas), pizzaRes)
	res = append(res, []string{fmt.Sprintf("%d", len(best.pizzas))})
	res = append(res, make([]string, len(best.pizzas)))
	for i := 0; i < len(best.pizzas); i++ {
		res[1][i] = fmt.Sprintf("%d", pizzaRes[i])
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

func (p *Proposition) IsBetterThan(p2 interface{}) bool {
	return p2 == nil || p2.(*Proposition) == nil || p.score > (p2.(*Proposition)).score
}
func (p *Proposition) IsMax() bool            { return p.score == max }
func (p *Proposition) CanHaveMoreChild() bool { return p.canHaveMoreChild }
func (p *Proposition) SetHeapIndex(index int) { p.heapIndex = index }
func (p *Proposition) GetHeapIndex() int      { return p.heapIndex }

func New(pizzas []int, parentSwapFrom int, baseScore int64) *Proposition {
	Proposition := &Proposition{pizzas, 0, baseScore, 0, parentSwapFrom, pizzas[parentSwapFrom] + 1, true}
	return Proposition
}

// Create n children from current Proposition
func (p *Proposition) CreateChildren() []interface{} {
	res := make([]interface{}, maxChildrenPerLoop)
	for i := 0; i < maxChildrenPerLoop; i++ {
		res[i] = NewChild(p)
	}
	return res
}

// Fill with pizzas until it throws up
func (p *Proposition) Compute() { // -> might be the most costly part of the algo > to be put in separate threads using batchExecution ?
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

		if i >= len(data) {
			println(i, len(data), size)
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

// This method is used to generate a new Proposition from an existing one
// As of now it is very basic meaning bad perfs
// -> given a Proposition, a new one will be attempted by switching one of the highest number of parts with the next lowest one
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

// Remove item at index lastSwapFrom and add item at the end with value lastSwapTo
func swapPizzas(p *Proposition) []int {
	newPizzas := make([]int, len(p.pizzas))
	copy(newPizzas[:p.lastSwapFrom], p.pizzas[:p.lastSwapFrom])
	copy(newPizzas[p.lastSwapFrom:], p.pizzas[p.lastSwapFrom+1:])
	newPizzas[len(newPizzas)-1] = p.lastSwapTo
	return newPizzas
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
