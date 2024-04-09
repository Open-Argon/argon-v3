package main

import (
	"fmt"
)

var threadCount = 0
var threadChan = make(chan bool)

func ArThread(args ...any) (any, ArErr) {
	if len(args) != 1 {
		return nil, ArErr{TYPE: "Type Error", message: "Invalid number of arguments, expected 1, got " + fmt.Sprint(len(args)), EXISTS: true}
	}
	var tocall any
	switch x := args[0].(type) {
	case anymap:
		if _, ok := x["__call__"]; !ok {
			return nil, ArErr{TYPE: "Type Error", message: "Cannot call thread with a '" + typeof(args[0]) + "'", EXISTS: true}
		}
		tocall = x["__call__"]
	case builtinFunc, Callable:
		tocall = x
	default:
		return nil, ArErr{TYPE: "Type Error", message: "Cannot call thread with a '" + typeof(args[0]) + "'", EXISTS: true}
	}
	var resp any
	var err ArErr

	hasrun := false
	joined := false
	var wg = make(chan bool)
	threadMap := Map(anymap{
		"start": builtinFunc{"start", func(args ...any) (any, ArErr) {
			if hasrun {
				return nil, ArErr{TYPE: "Runtime Error", message: "Cannot start a thread twice", EXISTS: true}
			}
			if len(args) != 0 {
				return nil, ArErr{TYPE: "Type Error", message: "Invalid number of arguments, expected 0, got " + fmt.Sprint(len(args)), EXISTS: true}
			}
			hasrun = true
			threadCount++
			go func() {
				resp, err = runCall(call{Callable: tocall, Args: []any{}}, nil, 0)
				wg <- true
				threadCount--
				if threadCount == 0 {
					threadChan <- true
				}
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
				return nil, ArErr{TYPE: "Type Error", message: "Invalid number of arguments, expected 0, got " + fmt.Sprint(len(args)), EXISTS: true}
			}
			joined = true
			<-wg
			return resp, err
		}},
	})
	return threadMap, ArErr{}
}
