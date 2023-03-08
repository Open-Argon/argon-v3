package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func anyToArgon(x any, quote bool, simplify bool, depth int, indent int) string {
	if depth == 0 {
		return "(...)"
	}
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
			if simplify {
				return numberToString(x, 0, true)
			}
			return numberToString(x, 0, false)
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
			output = append(output, anyToArgon(key, true, true, depth, indent+1)+": "+anyToArgon(x[key], true, true, depth-1, indent+1))
		}
		return "{\n" + (strings.Repeat("    ", indent+1)) + strings.Join(output, ",\n"+(strings.Repeat("    ", indent+1))) + "\n" + (strings.Repeat("    ", indent)) + "}"
	case builtinFunc:
		return "<builtin function " + x.name + ">"
	case Callable:
		return "<function " + x.name + ">"
	case ArClass:
		return anyToArgon(x.value, false, true, depth-1, indent+1)
	default:
		return fmt.Sprint(x)
	}
}
