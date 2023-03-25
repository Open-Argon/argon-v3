package main

import (
	"fmt"
	"sync"
)

func ArThread(args ...any) (any, ArErr) {
	if len(args) != 1 {
		return nil, ArErr{TYPE: "TypeError", message: "Invalid number of arguments, expected 1, got " + fmt.Sprint(len(args)), EXISTS: true}
	}
	var tocall any
	switch x := args[0].(type) {
	case anymap:
		if _, ok := x["__call__"]; !ok {
			return nil, ArErr{TYPE: "TypeError", message: "Cannot call thread with a '" + typeof(args[0]) + "'", EXISTS: true}
		}
		tocall = x["__call__"]
	case builtinFunc, Callable:
		tocall = x
	default:
		return nil, ArErr{TYPE: "TypeError", message: "Cannot call thread with a '" + typeof(args[0]) + "'", EXISTS: true}
	}
	var resp any
	var err ArErr

	hasrun := false
	joined := false
	var wg sync.WaitGroup
	threadMap := Map(anymap{
		"start": builtinFunc{"start", func(args ...any) (any, ArErr) {
			if hasrun {
				return nil, ArErr{TYPE: "Runtime Error", message: "Cannot start a thread twice", EXISTS: true}
			}
			if len(args) != 0 {
				return nil, ArErr{TYPE: "TypeError", message: "Invalid number of arguments, expected 0, got " + fmt.Sprint(len(args)), EXISTS: true}
			}
			hasrun = true
			wg.Add(1)
			go func() {
				resp, err = runCall(call{tocall, []any{}, "", 0, ""}, nil, 0)
				wg.Done()
			}()
			return nil, ArErr{}
		}},
		"join": builtinFunc{"join", func(args ...any) (any, ArErr) {
			if !hasrun {
				return nil, ArErr{TYPE: "Runtime Error", message: "Cannot join a thread that has not started", EXISTS: true}
			} else if joined {
				return nil, ArErr{TYPE: "Runtime Error", message: "Cannot join a thread twice", EXISTS: true}
			}
			if len(args) != 0 {
				return nil, ArErr{TYPE: "TypeError", message: "Invalid number of arguments, expected 0, got " + fmt.Sprint(len(args)), EXISTS: true}
			}
			joined = true
			wg.Wait()
			return resp, err
		}},
	})
	return threadMap, ArErr{}
}
