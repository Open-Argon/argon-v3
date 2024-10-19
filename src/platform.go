package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

var platform = Map(
	anymap{
		"os": ArString(runtime.GOOS),

		"cpu": Map(anymap{
			"count": newNumber().SetInt64((int64)(runtime.NumCPU())),
			"usage": builtinFunc{"usage", func(args ...any) (any, ArErr) {
				if len(args) != 2 {
					return nil, ArErr{
						TYPE:    "Runtime Error",
						message: "usage takes 2 arguments, got " + fmt.Sprint(len(args)),
						EXISTS:  true,
					}
				}
				if !isAnyNumber(args[0]) {
					return nil, ArErr{
						TYPE:    "Runtime Error",
						message: "first argument is meant to be a number, got " + fmt.Sprint(typeof(args[0])),
						EXISTS:  true,
					}
				}
				var Number = newNumber().Mul(args[0].(number), newNumber().SetInt64(1000)).Num().Int64()
				avgPercent, err := cpu.Percent(time.Duration(Number)*time.Millisecond, anyToBool(args[1]))
				if err != nil {
					return nil, ArErr{
						TYPE:    "Runtime Error",
						message: err.Error(),
						EXISTS:  true,
					}
				}
				var ArAvgPercent = []any{}
				for i := 0; i < len(avgPercent); i++ {
					ArAvgPercent = append(ArAvgPercent, newNumber().SetFloat64(avgPercent[i]))
				}
				return ArArray(ArAvgPercent), ArErr{}
			}},
		}),
	},
)
