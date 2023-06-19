package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/jwalton/go-supportscolor"
)

func anyToArgon(x any, quote bool, simplify bool, depth int, indent int, colored bool, plain int) string {
	x = ArValidToAny(x)
	output := []string{}
	maybenewline := ""
	if plain == 1 {
		maybenewline = "\n"
	}
	if colored {
		colored = supportscolor.Stdout().SupportsColor
	}
	if depth == 0 {
		if colored {
			output = append(output, color.New(38).Sprint("(...)"))
		} else {
			output = append(output, "(...)")
		}
		return strings.Join(output, "")
	}
	switch x := x.(type) {
	case string:
		if !quote {
			output = append(output, x)
			break
		}
		quoted := strconv.Quote(x)
		if colored {
			output = append(output, color.New(33).Sprint(quoted))
		} else {
			output = append(output, quoted)
		}
	case number:
		if colored {
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
			output = append(output, numberToString(x, simplify))
		}
		if colored {
			output = append(output, "\x1b[0m")
		}
	case bool:
		if colored {
			output = append(output, "\x1b[35;5;240m")
		}
		output = append(output, strconv.FormatBool(x))
		if colored {
			output = append(output, "\x1b[0m")
		}
	case nil:
		if colored {
			output = append(output, "\x1b[31;5;240m")
		}
		output = append(output, "null")
		if colored {
			output = append(output, "\x1b[0m")
		}
	case ArObject:

	case anymap:
		if len(x) == 0 {
			return "{}"
		}
		keys := make([]any, len(x))
		sort.Slice(keys, func(i, j int) bool {
			return anyToArgon(keys[i], false, true, 0, 0, false, 0) < anyToArgon(keys[j], false, true, 0, 0, false, 0)
		})

		i := 0
		for k := range x {
			keys[i] = k
			i++
		}
		output := []string{}
		for _, key := range keys {
			keyval := ""

			if typeof(key) != "string" || !SpacelessVariableCompiled.MatchString(key.(string)) {
				keyval = anyToArgon(key, true, true, depth, indent+1, colored, plain)
			} else {
				outputkeyval := []string{}
				if colored {
					outputkeyval = append(outputkeyval, "\x1b[36;5;240m")
				}
				outputkeyval = append(outputkeyval, key.(string))
				if colored {
					outputkeyval = append(outputkeyval, "\x1b[0m")
				}
				keyval = strings.Join(outputkeyval, "")
			}
			output = append(output, keyval+": "+anyToArgon(x[key], true, true, depth-1, indent+1, colored, plain))
		}
		return "{" + maybenewline + (strings.Repeat("    ", (indent+1)*plain)) + strings.Join(output, ","+maybenewline+(strings.Repeat("    ", (indent+1)*plain))) + maybenewline + (strings.Repeat("    ", indent*plain)) + "}"
	case []any:
		singleline := len(x) <= 3
		output := []string{}
		if simplify && len(x) >= 100 {
			for i := 0; i < 10; i++ {
				item := x[i]
				output = append(output, anyToArgon(item, true, true, depth-1, indent+1, colored, plain))
			}
			if colored {
				output = append(output, "\x1b[38;5;240m(...)\x1b[0m")
			} else {
				output = append(output, "(...)")
			}
			for i := len(x) - 10; i < len(x); i++ {
				item := x[i]
				output = append(output, anyToArgon(item, true, true, depth-1, indent+1, colored, plain))
			}
		} else {
			for i := 0; i < len(x); i++ {
				item := x[i]
				converted := anyToArgon(item, true, true, depth-1, indent+1, colored, plain)
				if singleline && strings.Contains(converted, "\n") {
					singleline = false
				}
				output = append(output, converted)
			}
		}

		if singleline {
			return "[" + strings.Join(output, ", ") + "]"
		}
		return "[" + maybenewline + (strings.Repeat("    ", (indent+1)*plain)) + strings.Join(output, ","+maybenewline+(strings.Repeat("    ", (indent+1)*plain))) + maybenewline + (strings.Repeat("    ", indent*plain)) + "]"
	case builtinFunc:
		if colored {
			output = append(output, "\x1b[38;5;240m")
		}
		output = append(output, "<builtin function "+x.name+">")
		if colored {
			output = append(output, "\x1b[0m")
		}
	case Callable:
		if colored {
			output = append(output, "\x1b[38;5;240m")
		}
		output = append(output, "<function "+x.name+">")
		if colored {
			output = append(output, "\x1b[0m")
		}
	default:
		return fmt.Sprint(x)
	}
	return strings.Join(output, "")
}
