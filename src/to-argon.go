package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func anyToArgon(x any, quote bool) string {
	switch x := x.(type) {
	case string:
		if !quote {
			return x
		}
		return strconv.Quote(x)
	case number:
		num, _ := x.Float64()
		if math.IsNaN(num) {
			return "NaN"
		} else if math.IsInf(num, 1) {
			return "infinity"
		} else if math.IsInf(num, -1) {
			return "-infinity"
		} else {
			return numberToString(x, 0)
		}
	case bool:
		return strconv.FormatBool(x)
	case nil:
		return "null"
	case ArMap:
		keys := make([]any, len(x))

		i := 0
		for k := range x {
			keys[i] = k
			i++
		}
		output := []string{}
		for _, key := range keys {
			output = append(output, anyToArgon(key, true)+": "+anyToArgon(x[key], true))
		}
		return "{" + strings.Join(output, ", ") + "}"
	case builtinFunc:
		return "<builtin function " + x.name + ">"
	case Callable:
		return "<function " + x.name + ">"
	case ArClass:
		return anyToArgon(x.value, false)
	default:
		return fmt.Sprint(x)
	}
}
