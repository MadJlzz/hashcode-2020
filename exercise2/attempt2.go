package exercise2

import (
	"fmt"
	"github.com/MadJlzz/hashcode-2020/algo/iterativeExecution"
	"github.com/MadJlzz/hashcode-2020/internal"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

var max float64
var data []float64
var size int
var sizeMinus1 int

const maxChildrenPerLoop = 100

func SolveExercise(fileContent map[int][]string) (res [][]string) {
	data = nil
	data = parseFile(fileContent)

	internal.SliceChildrenSize = len(data)
	internal.SliceChildrenCanAdd = func(a *internal.ArrayScored, i int) (bool, float64) {
		newScore := a.Score + data[i]
		return newScore <= max, newScore
	}

	start := &Proposition{&internal.SliceChildren{&internal.ArrayScored{}, &internal.ArrayScored{}, 0, true}, 0}

	best := iterativeExecution.IterativeExecution(start).(*Proposition)

	pizzaRes := reorder(best)
	fmt.Printf("score=%v length=%d res=%v\n", int64(best.s.Current.Score), len(best.s.Current.Data), pizzaRes)
	res = append(res, []string{fmt.Sprintf("%d", len(best.s.Current.Data))})
	res = append(res, make([]string, len(best.s.Current.Data)))
	for i := 0; i < len(best.s.Current.Data); i++ {
		res[1][i] = fmt.Sprintf("%d", pizzaRes[i])
	}
	return res
}

func reorder(p *Proposition) []int {
	var s = len(data) - 1
	res := make([]int, len(p.s.Current.Data))
	for i := 0; i < len(p.s.Current.Data); i++ {
		res[i] = s - p.s.Current.Data[i]
	}
	sort.Ints(res)
	return res
}

type Proposition struct {
	s         *internal.SliceChildren
	heapIndex int
}

func (p *Proposition) IsBetterThan(p2 interface{}) bool {
	return p2 == nil || p2.(*Proposition) == nil || p.s.Current.Score > (p2.(*Proposition)).s.Current.Score
}
func (p *Proposition) IsMax() bool            { return p.s.Current.Score == max }
func (p *Proposition) CanHaveMoreChild() bool { return p.s.CanHaveChild }
func (p *Proposition) SetHeapIndex(index int) { p.heapIndex = index }
func (p *Proposition) GetHeapIndex() int      { return p.heapIndex }
func (p *Proposition) Log() {
	fmt.Printf("%v %v %v\n", p.s.Current.Score, p.s.CanHaveChild, p.s.Current.Data)
}

// Create n children from current Proposition
func (p *Proposition) CreateChildren() []interface{} {
	res := make([]interface{}, maxChildrenPerLoop)
	for i := 0; i < maxChildrenPerLoop && p.s.CanHaveChild; i++ {
		temp := p.s.GetChild()
		if !reflect.ValueOf(temp).IsNil() {
			res[i] = &Proposition{temp, 0}
		}
	}
	return res
}

// Fill with pizzas until it throws up
func (p *Proposition) Compute() { // -> might be the most costly part of the algo > to be put in separate threads using batchExecution ?
	p.s.Complete()
}

func parseFile(fileContent map[int][]string) []float64 {
	size, _ = strconv.Atoi(strings.TrimSpace(fileContent[0][1]))
	res := make([]float64, size)
	sizeMinus1 = size - 1
	for i := 0; i < size; i++ {
		res[i], _ = strconv.ParseFloat(strings.TrimSpace(fileContent[1][sizeMinus1-i]), 64)
	}
	max, _ = strconv.ParseFloat(fileContent[0][0], 64)
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
