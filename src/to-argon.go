package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func anyToArgon(x any, quote bool, simplify bool, depth int, indent int, color bool, plain int) string {
	output := []string{}
	maybenewline := ""
	if plain == 1 {
		maybenewline = "\n"
	}
	if depth == 0 {
		if color {
			output = append(output, "\x1b[38;5;240m")
		}
		output = append(output, "(...)")
		if color {
			output = append(output, "\x1b[0m")
		}
		return strings.Join(output, "")
	}
	switch x := x.(type) {
	case string:
		if !quote {
			output = append(output, x)
			break
		}
		if color {
			output = append(output, "\x1b[33;5;240m")
		}
		output = append(output, strconv.Quote(x))
		if color {
			output = append(output, "\x1b[0m")
		}
	case number:
		if color {
			output = append(output, "\x1b[34;5;240m")
		}
		num, _ := x.Float64()
		if math.IsNaN(num) {
			output = append(output, "NaN")
		} else if math.IsInf(num, 1) {
			output = append(output, "infinity")
		} else if math.IsInf(num, -1) {
			output = append(output, "-infinity")
		} else {
			output = append(output, numberToString(x, 0, simplify))
		}
		if color {
			output = append(output, "\x1b[0m")
		}
	case bool:
		if color {
			output = append(output, "\x1b[35;5;240m")
		}
		output = append(output, strconv.FormatBool(x))
		if color {
			output = append(output, "\x1b[0m")
		}
	case nil:
		if color {
			output = append(output, "\x1b[31;5;240m")
		}
		output = append(output, "null")
		if color {
			output = append(output, "\x1b[0m")
		}
	case ArMap:
		if len(x) == 0 {
			return "{}"
		}
		keys := make([]any, len(x))

		i := 0
		for k := range x {
			keys[i] = k
			i++
		}
		output := []string{}
		for _, key := range keys {
			output = append(output, anyToArgon(key, true, true, depth, (indent+1)*plain, color, plain)+": "+anyToArgon(x[key], true, true, depth-1, indent+1, color, plain))
		}
		return "{" + maybenewline + (strings.Repeat("    ", (indent+1)*plain)) + strings.Join(output, ","+maybenewline+(strings.Repeat("    ", (indent+1)*plain))) + maybenewline + (strings.Repeat("    ", indent*plain)) + "}"
	case builtinFunc:
		if color {
			output = append(output, "\x1b[38;5;240m")
		}
		output = append(output, "<builtin function "+x.name+">")
		if color {
			output = append(output, "\x1b[0m")
		}
	case Callable:
		if color {
			output = append(output, "\x1b[38;5;240m")
		}
		output = append(output, "<function>")
		if color {
			output = append(output, "\x1b[0m")
		}
	case ArClass:
		return anyToArgon(x.value, quote, simplify, depth, indent, color, plain)
	default:
		return fmt.Sprint(x)
	}
	return strings.Join(output, "")
}
