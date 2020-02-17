package main

import (
	"flag"
	"github.com/madjlzz/hashcode-2020/exercise"
)

var (
	filename  = flag.String("filename", "test/e_also_big.in", "the data file used for our algorithm")
	algorithm = flag.String("solver", "", "the solver to use when trying to resolve the problem")
)

func main() {
	flag.Parse()

	exercise.ComputeFile(*filename)

	// 1. Reading input from file.
	//data := solver.ReadInput(*filename)

	//// 2. Transforming our data in a user friendly format.
	//
	//// 3. Retrieving the right solver and run the algorithm.
	//var s solver.Solver
	//switch *algorithm {
	//default:
	//	s = &solver.EmptySolver{}
	//}
	//
	//sol := s.Solve(data)
	//
	//// 4. Write that solution in a file.
	//
	//fmt.Printf("A solution has been found! %v\n", sol)
}
