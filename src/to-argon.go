package main

import (
	"fmt"
	"math"
	"strconv"
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
			return strconv.FormatFloat(num, 'f', -1, 64)
		}
	default:
		return fmt.Sprint(x)
	}
}
