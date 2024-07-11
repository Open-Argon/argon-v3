package main

import (
	"fmt"
	"math/rand"
	"time"
)

var rand_source = rand.New(rand.NewSource(time.Now().UnixMicro()))

func random() ArObject {
	return Number(rand_source.Float64())
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
	min := args[0].(ArObject)
	max := args[1].(ArObject)

	compare_num, err := CompareObjects(min, max)
	if err.EXISTS {
		return nil, err
	}

	compare, Err := numberToInt64(compare_num)
	if Err != nil {
		return nil, ArErr{
			TYPE:    "Runtime Error",
			message: Err.Error(),
			EXISTS:  true,
		}
	}

	if compare == 1 {
		return nil, ArErr{
			TYPE:    "Runtime Error",
			message: "range() num 1 must be less than or equal to num 2",
			EXISTS:  true,
		}
	}

	num_range, err := runOperation(
		operationType{
			operation: 11,
			values:    []any{max, min},
		},
		stack{},
		0,
	)

	if err.EXISTS {
		return nil, err
	}

	if _, ok := num_range.(ArObject); !ok {
		return nil, ArErr{
			TYPE:    "Runtime Error",
			message: "could not subtract the two numbers to calculate the range",
			EXISTS:  true,
		}
	}

	num_range_obj := num_range.(ArObject)

	rand := random()

	multiplier, err := runOperation(
		operationType{
			operation: 12,
			values:    []any{rand, num_range_obj},
		},
		stack{},
		0,
	)

	if err.EXISTS {
		return nil, err
	}

	if _, ok := multiplier.(ArObject); !ok {
		return nil, ArErr{
			TYPE:    "Runtime Error",
			message: "could not multiply the random number by the range",
			EXISTS:  true,
		}
	}

	return runOperation(
		operationType{
			operation: 10,
			values:    []any{multiplier, min},
		},
		stack{},
		0,
	)
}

var ArRandom = ArObject{anymap{
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
		new_seed, err := numberToInt64(a[0].(ArObject))
		if err != nil {
			return nil, ArErr{
				TYPE:    "Runtime Error",
				message: err.Error(),
				EXISTS:  true,
			}
		}
		rand_source.Seed(new_seed)
		return nil, ArErr{}
	}},
}}
