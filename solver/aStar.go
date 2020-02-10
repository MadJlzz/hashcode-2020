package solver

import (
	"container/heap"
	"fmt"
	"math"
)

// NEED TO BE IMPLEMENTED WITH POINTER
// ex: func (s *MyNode) SetFScore(score float64) {...}
type Noder interface {
	GetId() int
	GetNeighbours() []Noder
	GetDistance(neighbour Noder) float64
	SetFScore(score float64)
	GetFScore() float64

	BasicHeaper
}

// 'Expected' distance calculation method between 2 distant points
type HeuristicDistance func(node, goal Noder) float64

// Basic map but producing 'infinite' value if item is not present
type scoreMap map[int]float64

func (s *scoreMap) get(id int) float64 {
	if val, ok := (*s)[id]; ok {
		return val
	}
	return math.MaxFloat64
}

// Reconstruct the path computed in AStar to get the expected result
func getPath(cameFrom map[int]Noder, goal Noder) []Noder {
	current := goal
	var totalPath []Noder
	for ok := true; ok; current, ok = cameFrom[current.GetId()] {
		totalPath = append(totalPath, current)
	}

	reversePath := make([]Noder, len(totalPath))
	for i, j := len(totalPath)-1, 0; i >= 0; i, j = i-1, j+1 {
		reversePath[j] = totalPath[i]
	}
	return reversePath
}

// AStar algo as described in https://en.wikipedia.org/wiki/A*_search_algorithm
// staticWeight is used to increase algo speed at the expense of reliability -> 1.0 for perfect path, > 1.0 for less exact path but quicker result
func AStar(start Noder, goal Noder, distance HeuristicDistance, staticWeight float64) ([]Noder, error) {
	// Use of a map to efficiently retrieve an item
	openSet := map[int]Noder{start.GetId(): start}

	// Use of a heap to efficiently retrieve the item with the lowest score
	// NEED NODER.COMPARE TO BE LIKE this < A
	openSetHeap := &BasicHeap{start}
	heap.Init(openSetHeap)

	cameFrom := make(map[int]Noder)

	gScore := scoreMap{start.GetId(): 0}
	start.SetFScore(staticWeight * distance(start, goal))

	for len(openSet) > 0 {
		currentNode := heap.Pop(openSetHeap).(Noder)
		if currentNode.GetId() == goal.GetId() {
			return getPath(cameFrom, goal), nil
		}

		delete(openSet, currentNode.GetId())
		currentGScore := gScore.get(currentNode.GetId())
		for _, neighbour := range currentNode.GetNeighbours() {
			tentativeGScore := currentGScore + currentNode.GetDistance(neighbour)
			if tentativeGScore < gScore.get(neighbour.GetId()) {
				cameFrom[neighbour.GetId()] = currentNode
				gScore[neighbour.GetId()] = tentativeGScore
				neighbour.SetFScore(tentativeGScore + staticWeight*distance(neighbour, goal))

				if _, ok := openSet[neighbour.GetId()]; !ok {
					openSet[neighbour.GetId()] = neighbour
					heap.Push(openSetHeap, neighbour)
				} else {
					// As neighbour score was changed, heap needs to be updated
					heap.Fix(openSetHeap, neighbour.GetHeapIndex())
				}
			}
		}
	}
	return nil, fmt.Errorf("could not reach the goal")
}
