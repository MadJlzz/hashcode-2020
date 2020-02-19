package aStar

import (
	"fmt"
	"github.com/MadJlzz/hashcode-2020/tools"
	"math"
	"testing"
)

// Simple graph on 9 points forming a square will be used for this test.
// Starting point 1_1 is the one on bottom left
// Goal point 3_3 is the one on upper right
// Node can be linked to any direct neighbour (even diagonal)
// Distance is calculated using sqrt((x1-x2)² + (y1-y2)²)

func euclidianDist(na, nb Noder) float64 {
	a := na.(*testNode)
	b := nb.(*testNode)
	return math.Sqrt((a.x-b.x)*(a.x-b.x) + (a.y-b.y)*(a.y-b.y))
}

type testNode struct {
	id, heapIndex int
	fScore, x, y  float64
	neighbours    []Noder
}

func (t *testNode) GetId() int                           { return t.id }
func (t *testNode) GetNeighbours() []Noder               { return t.neighbours }
func (t *testNode) GetDistance(neighbour Noder) float64  { return euclidianDist(t, neighbour) }
func (t *testNode) SetFScore(score float64)              { t.fScore = score }
func (t *testNode) GetFScore() float64                   { return t.fScore }
func (t *testNode) IsBetterThan(heaper interface{}) bool { return t.fScore < heaper.(*testNode).fScore }
func (t *testNode) SetHeapIndex(index int)               { t.heapIndex = index }
func (t *testNode) GetHeapIndex() int                    { return t.heapIndex }

func newNode(id int, x, y float64) *testNode {
	node := &testNode{id: id, x: x, y: y}
	return node
}

var n1_1 = newNode(11, 1.0, 1.0)
var n1_2 = newNode(12, 1.0, 2.0)
var n1_3 = newNode(13, 1.0, 3.0)
var n2_1 = newNode(21, 2.0, 1.0)
var n2_2 = newNode(22, 2.0, 2.0)
var n2_3 = newNode(23, 2.0, 3.0)
var n3_1 = newNode(31, 3.0, 1.0)
var n3_2 = newNode(32, 3.0, 2.0)
var n3_3 = newNode(33, 3.0, 3.0)

func TestAStar_0(t *testing.T) {
	n1_1.neighbours = []Noder{n1_2, n2_1, n2_2}

	res, err := AStar(n1_1, n1_1, euclidianDist, 1.0)
	tools.UnittestAssert(t, nil, err)
	tools.UnittestAssert(t, 1, len(res))
	tools.UnittestAssert(t, 11, res[0].GetId())
}

/*
 * o -> node
 * s -> starting node
 * g -> goal node
 * y -> node part of the solution (shortest path)
 * x -> blocked path between 2 nodes
 *
 * o o g | o y g | oxo g | o oxo
 *       |    x  |  xxx  |    xx
 * o y o | o y o | o oxy | o o o
 *       |   x   |    x  |
 * s o o | s o o | s y y | o o o
 *       |       |       |
 *   1   |   2   |   3   |   4
 */
func TestAStar_1(t *testing.T) {
	n1_1.neighbours = []Noder{n1_2, n2_1, n2_2}
	n1_2.neighbours = []Noder{n1_1, n2_1, n2_2, n2_3, n1_3}
	n1_3.neighbours = []Noder{n1_2, n2_2, n2_3}
	n2_1.neighbours = []Noder{n3_1, n3_2, n2_2, n1_2, n1_1}
	n2_2.neighbours = []Noder{n2_1, n3_1, n3_2, n3_3, n2_3, n1_3, n1_2, n1_1}
	n2_3.neighbours = []Noder{n2_2, n3_2, n3_3, n1_3, n1_2}
	n3_1.neighbours = []Noder{n3_2, n2_2, n2_1}
	n3_2.neighbours = []Noder{n3_1, n3_3, n2_3, n2_2, n2_1}
	n3_3.neighbours = []Noder{n3_2, n2_3, n2_2}

	res, err := AStar(n1_1, n3_3, euclidianDist, 1.0)
	tools.UnittestAssert(t, nil, err)
	tools.UnittestAssert(t, 3, len(res))
	tools.UnittestAssert(t, 11, res[0].GetId())
	tools.UnittestAssert(t, 22, res[1].GetId())
	tools.UnittestAssert(t, 33, res[2].GetId())
}

/*
* o o g | o y g | oxo g | o oxo
*       |    x  |  xxx  |    xx
* o y o | o y o | o oxy | o o o
*       |   x   |    x  |
* s o o | s o o | s y y | o o o
*       |       |       |
*   1   |   2   |   3   |   4
 */
func TestAStar_2(t *testing.T) {
	n1_1.neighbours = []Noder{n1_2, n2_1, n2_2}
	n1_2.neighbours = []Noder{n1_1, n2_1, n2_3, n1_3}
	n1_3.neighbours = []Noder{n1_2, n2_2, n2_3}
	n2_1.neighbours = []Noder{n3_1, n3_2, n2_2, n1_2, n1_1}
	n2_2.neighbours = []Noder{n2_1, n3_1, n3_2, n2_3, n1_3, n1_1}
	n2_3.neighbours = []Noder{n2_2, n3_2, n3_3, n1_3, n1_2}
	n3_1.neighbours = []Noder{n3_2, n2_2, n2_1}
	n3_2.neighbours = []Noder{n3_1, n3_3, n2_3, n2_2, n2_1}
	n3_3.neighbours = []Noder{n3_2, n2_3}

	res, err := AStar(n1_1, n3_3, euclidianDist, 1.0)
	tools.UnittestAssert(t, nil, err)
	tools.UnittestAssert(t, 4, len(res))
	tools.UnittestAssert(t, 11, res[0].GetId())
	tools.UnittestAssert(t, 22, res[1].GetId())
	tools.UnittestAssert(t, 23, res[2].GetId())
	tools.UnittestAssert(t, 33, res[3].GetId())
}

/*
* o o g | o y g | oxo g | o oxo
*       |    x  |  xxx  |    xx
* o y o | o y o | o oxy | o o o
*       |   x   |    x  |
* s o o | s o o | s y y | o o o
*       |       |       |
*   1   |   2   |   3   |   4
 */
func TestAStar_3(t *testing.T) {
	n1_1.neighbours = []Noder{n1_2, n2_1, n2_2}
	n1_2.neighbours = []Noder{n1_1, n2_1, n2_2, n1_3}
	n1_3.neighbours = []Noder{n1_2, n2_3}
	n2_1.neighbours = []Noder{n3_1, n2_2, n1_2, n1_1}
	n2_2.neighbours = []Noder{n2_1, n1_2, n1_1}
	n2_3.neighbours = []Noder{n3_3, n1_3}
	n3_1.neighbours = []Noder{n2_1}
	n3_2.neighbours = []Noder{n3_3}
	n3_3.neighbours = []Noder{n3_2, n2_3}

	res, err := AStar(n1_1, n3_3, euclidianDist, 1.0)
	tools.UnittestAssert(t, nil, err)
	tools.UnittestAssert(t, 5, len(res))
	tools.UnittestAssert(t, 11, res[0].GetId())
	tools.UnittestAssert(t, 12, res[1].GetId())
	tools.UnittestAssert(t, 13, res[2].GetId())
	tools.UnittestAssert(t, 23, res[3].GetId())
	tools.UnittestAssert(t, 33, res[4].GetId())
}

/* o o g | o y g | oxo g | o oxo
*        |    x  |  xxx  |    xx
*  o y o | o y o | o oxy | o o o
*        |   x   |    x  |
*  s o o | s o o | s y y | o o o
*        |       |       |
*    1   |   2   |   3   |   4
 */
func TestAStar_4(t *testing.T) {
	n1_1.neighbours = []Noder{n1_2, n2_1, n2_2}
	n1_2.neighbours = []Noder{n1_1, n2_1, n2_2, n2_3, n1_3}
	n1_3.neighbours = []Noder{n1_2, n2_2, n2_3}
	n2_1.neighbours = []Noder{n3_1, n3_2, n2_2, n1_2, n1_1}
	n2_2.neighbours = []Noder{n2_1, n3_1, n3_2, n2_3, n1_3, n1_2, n1_1}
	n2_3.neighbours = []Noder{n2_2, n3_2, n1_3, n1_2}
	n3_1.neighbours = []Noder{n3_2, n2_2, n2_1}
	n3_2.neighbours = []Noder{n3_1, n2_3, n2_2, n2_1}
	n3_3.neighbours = []Noder{}

	res, err := AStar(n1_1, n3_3, euclidianDist, 1.0)
	tools.UnittestAssert(t, "could not reach the goal", fmt.Sprintf("%v", err))
	tools.UnittestAssert(t, 0, len(res))
}
