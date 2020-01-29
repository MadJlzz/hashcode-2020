package solver

import "log"

type Solver interface {
	Solve(data interface{}) interface{}
}

type EmptySolver struct{}

func (es *EmptySolver) Solve(data interface{}) interface{} {
	log.Fatalf("EmptySolver is here just as a placeholder for the solver skeleton! Please implement your own solver.")
	return nil
}
