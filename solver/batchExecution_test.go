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
	check(t, 0, res[0].Id)
	check(t, 1, res[1].Id)
	check(t, 2, res[2].Id)

	check(t, 2, res[0].Res)
	check(t, 6, res[1].Res)
	check(t, 8, res[2].Res)

	check(t, nil, res[0].Err)
	check(t, nil, res[1].Err)
	check(t, nil, res[2].Err)
}

func TestBatchExecutionError(t *testing.T) {
	data := []interface{}{DataTest{1}, DataTest{2}, DataTest{3}}
	res := BatchExecution(data, execute, 0)
	check(t, 3, len(res))
	check(t, 0, res[0].Id)
	check(t, 1, res[1].Id)
	check(t, 2, res[2].Id)

	check(t, 2, res[0].Res)
	check(t, nil, res[1].Res)
	check(t, 6, res[2].Res)

	check(t, nil, res[0].Err)
	check(t, "2 is special", fmt.Sprintf("%s", res[1].Err))
	check(t, nil, res[2].Err)
}

func TestBatchExecutionErrorWithTimeout(t *testing.T) {
	data := []interface{}{DataTest{1}, DataTest{2}, DataTest{42}, DataTest{42}, DataTest{3}}
	res := BatchExecution(data, execute, 100) // test should work if your computer is not from the 90s
	check(t, 5, len(res))
	check(t, 0, res[0].Id)
	check(t, 1, res[1].Id)
	check(t, 2, res[2].Id)
	check(t, 3, res[3].Id)
	check(t, 4, res[4].Id)

	check(t, 2, res[0].Res)
	check(t, nil, res[1].Res)
	check(t, nil, res[2].Res)
	check(t, nil, res[3].Res)
	check(t, 6, res[4].Res)

	check(t, nil, res[0].Err)
	check(t, "2 is special", fmt.Sprintf("%s", res[1].Err))
	check(t, "timeout reached", fmt.Sprintf("%s", res[2].Err))
	check(t, "timeout reached", fmt.Sprintf("%s", res[3].Err))
	check(t, nil, res[4].Err)
}

func check(t *testing.T, want interface{}, got interface{}) {
	if want != got {
		t.Errorf("Error: expected %v - got %v", want, got)
	}
}
