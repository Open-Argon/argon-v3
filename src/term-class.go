package main

import (
	"fmt"
	"time"
)

var timing = ArMap{}

var plain = ArMap{
	"log": builtinFunc{"log", func(args ...any) (any, ArErr) {
		output := []any{}
		for i := 0; i < len(args); i++ {
			output = append(output, anyToArgon(args[i], false, true, 3, 0, false, 0))
		}
		fmt.Println(output...)
		return nil, ArErr{}
	}},
	"logVal": builtinFunc{"logVal", func(args ...any) (any, ArErr) {
		output := []any{}
		for i := 0; i < len(args); i++ {
			output = append(output, anyToArgon(args[i], true, true, 3, 0, false, 0))
		}
		fmt.Println(output...)
		return nil, ArErr{}
	}},
}

var ArTerm = ArMap{
	"log": builtinFunc{"log", func(args ...any) (any, ArErr) {
		output := []any{}
		for i := 0; i < len(args); i++ {
			output = append(output, anyToArgon(args[i], false, true, 3, 0, true, 1))
		}
		fmt.Println(output...)
		return nil, ArErr{}
	}},
	"logVal": builtinFunc{"logVal", func(args ...any) (any, ArErr) {
		output := []any{}
		for i := 0; i < len(args); i++ {
			output = append(output, anyToArgon(args[i], true, true, 3, 0, true, 1))
		}
		fmt.Println(output...)
		return nil, ArErr{}
	}},
	"print": builtinFunc{"print", func(args ...any) (any, ArErr) {
		output := []any{}
		for i := 0; i < len(args); i++ {
			output = append(output, anyToArgon(args[i], false, false, 3, 0, false, 1))
		}
		fmt.Println(output...)
		return nil, ArErr{}
	}},
	"plain": plain,
	"error": builtinFunc{"error", func(args ...any) (any, ArErr) {
		output := []any{"error: "}
		for i := 0; i < len(args); i++ {
			output = append(output, anyToArgon(args[i], false, true, 3, 0, false, 0))
		}
		fmt.Printf("\x1b[%dm%s\x1b[0m", 91, fmt.Sprint(output...)+"\n")
		return nil, ArErr{}
	},
	},
	"warn": builtinFunc{"error", func(args ...any) (any, ArErr) {
		output := []any{"warning: "}
		for i := 0; i < len(args); i++ {
			output = append(output, anyToArgon(args[i], false, true, 3, 0, false, 0))
		}
		fmt.Printf("\x1b[%dm%s\x1b[0m", 93, fmt.Sprint(output...)+"\n")
		return nil, ArErr{}
	},
	},
	"time": builtinFunc{"time", func(args ...any) (any, ArErr) {
		var id any = nil
		if len(args) > 0 {
			id = args[0]
		}
		timing[id] = time.Now()
		return nil, ArErr{}
	},
	},
	"timeEnd": builtinFunc{"timeEnd", func(args ...any) (any, ArErr) {
		var id any = nil
		if len(args) > 0 {
			id = args[0]
		}
		if _, ok := timing[id]; !ok {
			return nil, ArErr{TYPE: "TypeError", message: "Cannot find timer with id '" + fmt.Sprint(id) + "'", EXISTS: true}
		}
		timesince := time.Since(timing[id].(time.Time))
		delete(timing, id)
		if id == nil {
			id = "Timer"
		}
		fmt.Printf("\x1b[%dm%s\x1b[0m", 34, fmt.Sprint(anyToArgon(id, false, true, 3, 0, false, 0), ": ", timesince)+"\n")
		return nil, ArErr{}
	}},
}
