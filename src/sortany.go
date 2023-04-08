package main

import "fmt"

type keyCache map[any]any

func quickSort(list []any, getKey func(any) (any, ArErr)) ([]any, ArErr) {
	if len(list) <= 1 {
		return list, ArErr{}
	}

	pivot := list[0]
	var left []any
	var right []any

	var cache = make(keyCache)

	for _, v := range list[1:] {
		val, err := getkeyCache(getKey, v, cache)
		if err.EXISTS {
			return nil, err
		}
		pivotval, err := getkeyCache(getKey, pivot, cache)
		if err.EXISTS {
			return nil, err
		}
		comp, comperr := compare(val, pivotval)
		if comperr != nil {
			return nil, ArErr{
				TYPE:    "TypeError",
				message: comperr.Error(),
				EXISTS:  true,
			}
		}
		if comp {
			left = append(left, v)
		} else {
			right = append(right, v)
		}
	}

	left, err := quickSort(left, getKey)
	if err.EXISTS {
		return nil, err
	}
	right, err = quickSort(right, getKey)
	if err.EXISTS {
		return nil, err
	}

	return append(append(left, pivot), right...), ArErr{}
}

func getkeyCache(getKey func(any) (any, ArErr), index any, cache keyCache) (any, ArErr) {
	key := ArValidToAny(index)
	if cacheval, ok := cache[key]; ok {
		return cacheval, ArErr{}
	}
	val, err := getKey(index)
	if err.EXISTS {
		return nil, err
	}
	cache[key] = val
	return val, ArErr{}
}

func compare(a, b any) (bool, error) {
	if isAnyNumber(a) && isAnyNumber(b) {
		return a.(number).Cmp(b.(number)) < 0, nil
	} else if x, ok := a.(ArObject); ok {
		if y, ok := x.obj["__LessThan__"]; ok {
			resp, err := runCall(
				call{
					callable: y,
					args:     []any{b},
				}, stack{}, 0)
			if !err.EXISTS {
				return anyToBool(resp), nil
			}
		}
	} else if x, ok := b.(byte); ok {
		if y, ok := a.(byte); ok {
			return y < x, nil
		}
	}
	return false, fmt.Errorf("cannot compare %s to %s", typeof(a), typeof(b))
}
