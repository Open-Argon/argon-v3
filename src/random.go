package main

import (
	"fmt"
	"math/rand"
	"time"
)

func random() number {
	return newNumber().SetFloat64(
		rand.Float64(),
	)
}

func randomRange(args ...any) (any, ArErr) {
	if len(args) != 2 {
		return nil, ArErr{
			TYPE:    "Runtime Error",
			message: "takes 2 arguments, got " + fmt.Sprint(len(args)),
			EXISTS:  true,
		}
	}
	if typeof(args[0]) != "number" {
		return nil, ArErr{
			TYPE:    "Runtime Error",
			message: "takes a number not a '" + typeof(args[0]) + "'",
			EXISTS:  true,
		}
	} else if typeof(args[1]) != "number" {
		return nil, ArErr{
			TYPE:    "Runtime Error",
			message: "takes a number not a '" + typeof(args[1]) + "'",
			EXISTS:  true,
		}
	}
	min := args[0].(number)
	max := args[1].(number)
	if min.Cmp(max) > 0 {
		return nil, ArErr{
			TYPE:    "Runtime Error",
			message: "takes a min less than max",
			EXISTS:  true,
		}
	}
	difference := newNumber().Sub(max, min)
	rand := random()
	rand.Mul(rand, difference)
	rand.Add(rand, min)
	return rand, ArErr{}
}

var ArRandom = Map(anymap{
	"__call__": builtinFunc{"random", func(args ...any) (any, ArErr) {
		if len(args) != 0 {
			return nil, ArErr{
				TYPE:    "Runtime Error",
				message: "takes 0 arguments, got " + fmt.Sprint(len(args)),
				EXISTS:  true,
			}
		}
		return random(), ArErr{}
	}},
	"int": builtinFunc{"int", func(a ...any) (any, ArErr) {
		resp, err := randomRange(a...)
		if err.EXISTS {
			return nil, err
		}
		return round(resp.(number), 0), ArErr{}
	}},
	"range": builtinFunc{"range", randomRange},
	"seed": builtinFunc{"seed", func(a ...any) (any, ArErr) {
		if len(a) != 1 {
			return nil, ArErr{
				TYPE:    "Runtime Error",
				message: "takes 1 argument, got " + fmt.Sprint(len(a)),
				EXISTS:  true,
			}
		}
		if typeof(a[0]) != "number" {
			return nil, ArErr{
				TYPE:    "Runtime Error",
				message: "takes a number not a '" + typeof(a[0]) + "'",
				EXISTS:  true,
			}
		}
		if !a[0].(number).IsInt() {
			return nil, ArErr{
				TYPE:    "Runtime Error",
				message: "takes an integer not a float",
				EXISTS:  true,
			}
		}
		rand.Seed(
			a[0].(number).Num().Int64(),
		)
		return nil, ArErr{}
	}},
})

func init() {
	rand.Seed(
		time.Now().UnixMicro(),
	)
}
