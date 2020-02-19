package genetic

import (
	"math"
	"testing"
	"time"
)

var ChildNb = 10000

type simple struct {
	squareRootOfTwo float64
	res             float64
	mother          *simple
}

func (s simple) GetMother() Celler { return s.mother }
func (s simple) GetDistance(celler Celler) []float64 {
	return []float64{math.Abs(s.squareRootOfTwo - s.mother.squareRootOfTwo)}
}
func (s simple) CreateMutants(iteration int, distance []float64) (res []Celler) {
	mutants := GenericMutateFloat(ChildNb, distance[0], s.squareRootOfTwo, -100, 100)
	res = make([]Celler, ChildNb)
	for i, v := range mutants {
		res[i] = simple{squareRootOfTwo: v, mother: &s}
	}
	return res
}

func (s simple) Score() float64 {
	score := math.Abs(s.res - 2)
	if score == 0 {
		return math.MaxFloat64
	}
	return 1 / (score * score)
}

func executeSimple(c Celler) Celler {
	s := c.(simple)
	s.res = s.squareRootOfTwo * s.squareRootOfTwo
	return s
}

func TestSimple(t *testing.T) {
	first := simple{squareRootOfTwo: 50, mother: &simple{squareRootOfTwo: 75, mother: nil}}
	res := AppliedDarwinism(first, executeSimple, 10, 1000)

	s := res.(simple)
	println(s.squareRootOfTwo, s.res)
	if math.Abs(s.res-2) > 0.0001 {
		t.Errorf("Not efficient enough")
	}
}

func TestApproximativeSimple(t *testing.T) {
	start := time.Now()
	first := simple{squareRootOfTwo: 50, mother: &simple{squareRootOfTwo: 75, mother: nil}}
	res := AppliedDarwinism(first, executeSimple, 100, 100)

	s := res.(simple)
	println(s.squareRootOfTwo, s.res)
	if math.Abs(s.res-2) > 0.1 {
		t.Errorf("Not efficient enough")
	}
	if time.Since(start) > 200*time.Millisecond {
		t.Errorf("Should have stopped by now")
	}
}

type slightlyComplex struct {
	squareRootOfTwo float64
	pie             float64
	fortyTwo        int

	res    float64
	mother *slightlyComplex
}

func (s slightlyComplex) GetMother() Celler { return s.mother }
func (s slightlyComplex) GetDistance(celler Celler) []float64 {
	return []float64{
		math.Abs(s.squareRootOfTwo - s.mother.squareRootOfTwo),
		math.Abs(s.pie - s.mother.pie),
		math.Abs(float64(s.fortyTwo) - float64(s.mother.fortyTwo)),
	}
}
func (s slightlyComplex) CreateMutants(iteration int, distance []float64) (res []Celler) {
	mutantsSQ := GenericMutateFloat(ChildNb, distance[0], s.squareRootOfTwo, 0, 100)
	mutantsPie := GenericMutateFloat(ChildNb, distance[1], s.pie, -100, 100)
	mutantsFT := GenericMutateInt(ChildNb, distance[2], s.fortyTwo, -100, 100)
	res = make([]Celler, ChildNb)
	for i, _ := range mutantsSQ {
		res[i] = slightlyComplex{squareRootOfTwo: mutantsSQ[i], pie: mutantsPie[i], fortyTwo: mutantsFT[i], mother: &s}
	}
	return res
}

func (s slightlyComplex) Score() float64 {
	if s.res == 0 {
		return math.MaxFloat64
	}
	return 1 / s.res
}

func executeSlightlyComplex(c Celler) Celler {
	s := c.(slightlyComplex)
	s.res = sqrt(2, s.squareRootOfTwo*s.squareRootOfTwo) + sqrt(42, float64(s.fortyTwo)) + sqrt(s.pie, math.Pi)
	return s
}

func sqrt(a float64, b float64) float64 {
	res := a - b
	return res * res
}

func TestSlightlyComplex(t *testing.T) {
	first := slightlyComplex{squareRootOfTwo: 50, pie: 50, fortyTwo: 50, mother: &slightlyComplex{squareRootOfTwo: 75, pie: 75, fortyTwo: 75, mother: nil}}
	res := AppliedDarwinism(first, executeSlightlyComplex, 10, 1000)

	s := res.(slightlyComplex)
	println(s.squareRootOfTwo, s.pie, s.fortyTwo, s.res)
	if math.Abs(s.res) > 0.00001 {
		t.Errorf("Not efficient enough")
	}
}
