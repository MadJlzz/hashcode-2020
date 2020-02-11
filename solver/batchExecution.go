package solver

import (
	"fmt"
	"time"
)

const BatchTimeout = 3600000

/*
Purpose is to provide the execution of N different jobs in parallel with a time limit if provided
All results are then sent back to be handled by the caller
*/

type ExecutionRes struct {
	Res interface{}
	Err error
}

type BatchRes struct {
	Id int
	ExecutionRes
}

// provide an easier to use batch execution interface. In case of error (technical or timeout), item will be null -> must be anticipated by the user
func BatchExecutionBasic(data []interface{}, execute func(interface{}) interface{}, timeout time.Duration) []interface{} {
	tempFunc := func(i interface{}) ExecutionRes {
		return ExecutionRes{execute(i), nil}
	}

	resRaw := BatchExecution(data, tempFunc, timeout)

	res := make([]interface{}, len(resRaw))
	for i, v := range resRaw {
		if v.Err != nil {
			res[i] = nil
		} else {
			res[i] = v.Res
		}
	}

	return res
}

func BatchExecution(data []interface{}, execute func(interface{}) ExecutionRes, timeout time.Duration) []*BatchRes {
	if timeout == 0 {
		// 0 means no timeout -> artificially set sufficiently large one -> 1h
		timeout = BatchTimeout
	}

	res := make([]*BatchRes, len(data))
	channel := make(chan *BatchRes)
	for id, val := range data {
		go executeSubroutine(id, channel, val, execute)
	}

	for i := 0; i < len(data); i++ {
		select {
		case tempRes := <-channel:
			res[tempRes.Id] = tempRes
		case <-time.After(timeout * time.Millisecond):
			return fillWithError(res)
		}
	}
	return res
}

func executeSubroutine(id int, channel chan<- *BatchRes, data interface{}, execute func(interface{}) ExecutionRes) {
	executionRes := execute(data)
	channel <- &BatchRes{Id: id, ExecutionRes: executionRes}
}

func fillWithError(batchRes []*BatchRes) []*BatchRes {
	for i, _ := range batchRes {
		if batchRes[i] == nil {
			batchRes[i] = &BatchRes{
				Id:           i,
				ExecutionRes: ExecutionRes{Err: fmt.Errorf("timeout reached")},
			}
		}
	}
	return batchRes
}
