package main

import "sync"

func ArThread(args ...any) (any, ArErr) {
	if len(args) == 0 {
		return nil, ArErr{TYPE: "TypeError", message: "Cannot call thread without a function", EXISTS: true}
	}
	var tocall any
	switch x := args[0].(type) {
	case Callable:
		tocall = x
	case builtinFunc:
		tocall = x
	default:
		return nil, ArErr{TYPE: "TypeError", message: "Cannot call thread with a '" + typeof(args[0]) + "'", EXISTS: true}
	}
	var resp any
	var err ArErr
	currentscope := stack{vars, scope{}}
	hasrun := false
	joined := false
	var wg sync.WaitGroup
	threaMap := ArMap{
		"start": builtinFunc{"start", func(args ...any) (any, ArErr) {
			if hasrun {
				return nil, ArErr{TYPE: "Runtime Error", message: "Cannot start a thread twice", EXISTS: true}
			}
			hasrun = true
			wg.Add(1)
			go func() {
				resp, err = runCall(call{tocall, []any{}, "", 0, ""}, currentscope, 0)
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
			joined = true
			wg.Wait()
			return resp, err
		}},
	}
	return threaMap, ArErr{}
}
