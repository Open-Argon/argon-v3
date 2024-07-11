package main

import "math/big"

func AnyToArValid(arr any) any {
	switch arr := arr.(type) {
	case []any:
		return ArArray(arr)
	case string:
		return ArString(arr)
	case anymap:
		return Map(arr)
	case []byte:
		return ArBuffer(arr)
	case byte:
		return ArByte(arr)
	case int, int64, float64, float32, *big.Rat, *big.Int:
		return Number(arr)
	default:
		return arr
	}
}

func ArValidToAny(a any) any {
	switch a := a.(type) {
	case ArObject:
		if v, ok := a.obj["__value__"]; ok {
			return v
		}
	}
	return a
}
