package main

import (
	"fmt"
	"time"
)

var timing = anymap{}

var ArTerm = Map(anymap{
	"log": builtinFunc{"log", func(args ...any) (any, ArErr) {
		output := []any{}
		for i := 0; i < len(args); i++ {
			output = append(output, anyToArgon(args[i], false, true, 3, 0, true, 1))
		}
		fmt.Println(output...)
		return nil, ArErr{}
	}},
	"clear": builtinFunc{"clear", func(args ...any) (any, ArErr) {
		if len(args) != 0 {
			return nil, ArErr{
				TYPE:    "Runtime Error",
				message: "takes 0 arguments, got " + fmt.Sprint(len(args)),
				EXISTS:  true,
			}
		}
		fmt.Print("\033[H\033[2J")
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
	"plain": Map(anymap{
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
		"print": builtinFunc{"print", func(args ...any) (any, ArErr) {
			output := []any{}
			for i := 0; i < len(args); i++ {
				output = append(output, anyToArgon(args[i], false, false, 3, 0, false, 0))
			}
			fmt.Println(output...)
			return nil, ArErr{}
		}},
		"oneLine": builtinFunc{"oneLine", func(args ...any) (any, ArErr) {
			output := []any{}
			for i := 0; i < len(args); i++ {
				output = append(output, anyToArgon(args[i], false, false, 3, 0, false, 0))
			}
			fmt.Print(output...)
			return nil, ArErr{}
		}},
	}),
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
			id = ArValidToAny(args[0])
		}
		timing[id] = time.Now()
		return nil, ArErr{}
	},
	},
	"timeEnd": builtinFunc{"timeEnd", func(args ...any) (any, ArErr) {
		var id any = nil
		if len(args) > 0 {
			id = ArValidToAny(args[0])
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
})

var ArInput = Map(
	anymap{
		"password": builtinFunc{"password", func(args ...any) (any, ArErr) {
			resp, err := getPassword(args...)
			if err != nil {
				return nil, ArErr{TYPE: "Runtime Error", message: err.Error(), EXISTS: true}
			}
			return ArString(resp), ArErr{}
		}},
		"__call__": builtinFunc{"input", func(args ...any) (any, ArErr) {
			return input(args...), ArErr{}
		}},
		"pause": builtinFunc{"pause", func(args ...any) (any, ArErr) {
			pause()
			return nil, ArErr{}
		}},
	},
)
