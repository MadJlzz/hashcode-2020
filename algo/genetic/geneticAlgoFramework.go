package genetic

import (
	"fmt"
	"github.com/MadJlzz/hashcode-2020/algo/batchExecution"
	"gonum.org/v1/gonum/stat/distuv"
	"math"
	"time"
)

const GeneticBaseDistance = 1
const GeneticBaseDistanceIncrease = 5

type Celler interface {
	// Create N cells from the current one with properties varying around their mother. Current is added to the result
	// Might be a good idea to keep track of distance between current cell and its mother to estimate the spread of its childs
	// i.e. -> if current was far from its mother, might be a good idea to keep a wide search -> have a good spread of properties between its children to have a wide search range
	//      -> if current was close to its mother, we are likely to get close to the optimal solution -> restrain the spread of properties to get as close to the optimal result as possible
	CreateMutants(iteration int, distance []float64) []Celler

	// Provide a global score for the current cell to help select the best one for the next iteration
	Score() float64

	// used to define distance to previous, to narrow the search (first cell must reference itself as mother -> the egg is the chicken)
	GetMother() Celler
	GetDistance(celler Celler) []float64
}

func AppliedDarwinism(celler Celler, execute func(Celler) Celler, maxIterations int, timeout time.Duration) Celler {
	start := time.Now()
	distances := celler.GetDistance(celler.GetMother()) // just to get the number of expected distances
	channel := make(chan Celler)

	bestCell := celler
	for i := 0; i < maxIterations; i++ {
		cells := bestCell.CreateMutants(0, distances)
		go geneticAlgoExecutions(channel, cells, execute, timeout)

		var temp Celler
		select {
		case temp = <-channel:
		case <-time.After(timeout*time.Millisecond - time.Since(start)):
			fmt.Printf("Timeout reached after %d iterations. Returning best result", i)
			return bestCell
		}

		if temp != nil {
			bestCell = temp
			distances = getNewDistances(distances, bestCell.GetDistance(bestCell.GetMother()))
		} else {
			println("No good cell found this turn ! Keeping last one")
			distances = increaseDistances(distances)
		}
	}
	return bestCell
}

// New Standard Deviation must be above 0. If not, get previous one but multiply it to increase search range. If no previous one, default to 1
func getNewDistances(oldDistances []float64, newDistances []float64) (res []float64) {
	res = make([]float64, len(newDistances))
	for i, _ := range newDistances {
		switch {
		case newDistances[i] != 0:
			res[i] = newDistances[i]
		case oldDistances[i] == 0:
			res[i] = GeneticBaseDistance
		default:
			res[i] = oldDistances[i] * GeneticBaseDistanceIncrease
		}
	}
	return res
}

func increaseDistances(oldDistances []float64) (res []float64) {
	res = make([]float64, len(oldDistances))
	for i, _ := range oldDistances {
		switch oldDistances[i] {
		case 0:
			res[i] = GeneticBaseDistance
		default:
			res[i] = oldDistances[i] * GeneticBaseDistanceIncrease
		}
	}
	return res
}

func geneticAlgoExecutions(channel chan<- Celler, cellers []Celler, execute func(Celler) Celler, timeout time.Duration) {
	batchFunc := func(cellI interface{}) batchExecution.ExecutionRes {
		cell := cellI.(Celler)
		res := execute(cell)
		if res == nil {
			return batchExecution.ExecutionRes{nil, fmt.Errorf("no result for cell %v", cell)}
		}
		return batchExecution.ExecutionRes{res, nil}
	}

	res := batchExecution.BatchExecution(toInterface(cellers), batchFunc, timeout/4)
	var bestRes Celler
	for _, v := range res {
		temp := v.Res.(Celler)
		if bestRes == nil || bestRes.Score() < temp.Score() {
			bestRes = temp
		}
	}
	channel <- bestRes
}

func toInterface(cellers []Celler) []interface{} {
	res := make([]interface{}, len(cellers))
	for i, v := range cellers {
		res[i] = v
	}
	return res
}

// returns N numbers distributed around start using standard deviation from input
// Assumes you're a good guy and not setting impossible conditions !
func GenericMutateFloat(childNb int, distance float64, start float64, min float64, max float64) (mutants []float64) {
	mutants = make([]float64, childNb)
	dist := distuv.Normal{
		Mu:    start,
		Sigma: math.Min(distance, max-min),
	}

	for i := 0; i < childNb; i++ {
		var rnd float64
		for rnd = dist.Rand(); rnd > max || rnd < min; rnd = dist.Rand() {
		}
		mutants[i] = rnd
	}
	return mutants
}

// returns N numbers distributed around start using standard deviation from input
// Assumes you're a good guy and not setting impossible conditions !
func GenericMutateInt(childNb int, distance float64, start int, min int, max int) (mutants []int) {
	mutants = make([]int, childNb)
	dist := distuv.Normal{
		Mu:    float64(start),
		Sigma: math.Min(distance, float64(max-min)),
	}

	getInt := func() int { return int(math.Round(dist.Rand())) }

	for i := 0; i < childNb; i++ {
		var rnd int
		for rnd = getInt(); rnd > max || rnd < min; rnd = getInt() {
		}
		mutants[i] = rnd
	}
	return mutants
}
