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
	res interface{}
	err error
}

type BatchRes struct {
	id int
	ExecutionRes
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
			res[tempRes.id] = tempRes
		case <-time.After(timeout * time.Millisecond):
			return fillWithError(res)
		}
	}
	return res
}

func executeSubroutine(id int, channel chan<- *BatchRes, data interface{}, execute func(interface{}) ExecutionRes) {
	executionRes := execute(data)
	channel <- &BatchRes{id: id, ExecutionRes: executionRes}
}

func fillWithError(batchRes []*BatchRes) []*BatchRes {
	for i, _ := range batchRes {
		if batchRes[i] == nil {
			batchRes[i] = &BatchRes{
				id:           i,
				ExecutionRes: ExecutionRes{err: fmt.Errorf("timeout reached")},
			}
		}
	}
	return batchRes
}
