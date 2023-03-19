package main

import "fmt"

type keyCache map[interface{}]interface{}

func quickSort(list []interface{}, getKey func(interface{}) (interface{}, ArErr)) ([]interface{}, ArErr) {
	if len(list) <= 1 {
		return list, ArErr{}
	}

	pivot := list[0]
	var left []interface{}
	var right []interface{}

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
		if comp < 0 {
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

func getkeyCache(getKey func(interface{}) (interface{}, ArErr), index interface{}, cache keyCache) (interface{}, ArErr) {
	if cacheval, ok := cache[index]; ok {
		return cacheval, ArErr{}
	}
	val, err := getKey(index)
	if err.EXISTS {
		return nil, err
	}
	cache[index] = val
	return val, ArErr{}
}

func compare(a, b interface{}) (int, error) {
	switch x := a.(type) {
	case string:
		if _, ok := b.(string); !ok {
			return 0, fmt.Errorf("cannot compare %T to %T", a, b)
		}
		if a == b {
			return 0, nil
		}
		if x < b.(string) {
			return -1, nil
		}
		return 1, nil
	case number:
		if _, ok := b.(number); !ok {
			return 0, fmt.Errorf("cannot compare %T to %T", a, b)
		}
		return x.Cmp(b.(number)), nil
	}
	return 0, fmt.Errorf("cannot compare %T to %T", a, b)
}
