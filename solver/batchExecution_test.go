package solver

import (
	"fmt"
	"testing"
	"time"
)

type DataTest struct {
	val int
}

func execute(dataTest interface{}) ExecutionRes {
	data := dataTest.(DataTest)
	switch data.val {
	case 42:
		time.Sleep(1 * time.Second) // nobody can compute 42 that quickly
		return ExecutionRes{42, nil}
	case 2:
		return ExecutionRes{nil, fmt.Errorf("2 is special")}
	default:
		return ExecutionRes{2 * data.val, nil}
	}
}

func TestBatchExecutionNominal(t *testing.T) {
	data := []interface{}{DataTest{1}, DataTest{3}, DataTest{4}}
	res := BatchExecution(data, execute, 0)
	check(t, 3, len(res))
	check(t, 0, res[0].id)
	check(t, 1, res[1].id)
	check(t, 2, res[2].id)

	check(t, 2, res[0].res)
	check(t, 6, res[1].res)
	check(t, 8, res[2].res)

	check(t, nil, res[0].err)
	check(t, nil, res[1].err)
	check(t, nil, res[2].err)
}

func TestBatchExecutionError(t *testing.T) {
	data := []interface{}{DataTest{1}, DataTest{2}, DataTest{3}}
	res := BatchExecution(data, execute, 0)
	check(t, 3, len(res))
	check(t, 0, res[0].id)
	check(t, 1, res[1].id)
	check(t, 2, res[2].id)

	check(t, 2, res[0].res)
	check(t, nil, res[1].res)
	check(t, 6, res[2].res)

	check(t, nil, res[0].err)
	check(t, "2 is special", fmt.Sprintf("%s", res[1].err))
	check(t, nil, res[2].err)
}

func TestBatchExecutionErrorWithTimeout(t *testing.T) {
	data := []interface{}{DataTest{1}, DataTest{2}, DataTest{42}, DataTest{42}, DataTest{3}}
	res := BatchExecution(data, execute, 100) // test should work if your computer is not from the 90s
	check(t, 5, len(res))
	check(t, 0, res[0].id)
	check(t, 1, res[1].id)
	check(t, 2, res[2].id)
	check(t, 3, res[3].id)
	check(t, 4, res[4].id)

	check(t, 2, res[0].res)
	check(t, nil, res[1].res)
	check(t, nil, res[2].res)
	check(t, nil, res[3].res)
	check(t, 6, res[4].res)

	check(t, nil, res[0].err)
	check(t, "2 is special", fmt.Sprintf("%s", res[1].err))
	check(t, "timeout reached", fmt.Sprintf("%s", res[2].err))
	check(t, "timeout reached", fmt.Sprintf("%s", res[3].err))
	check(t, nil, res[4].err)
}

func check(t *testing.T, want interface{}, got interface{}) {
	if want != got {
		t.Errorf("Error: expected %v - got %v", want, got)
	}
}
