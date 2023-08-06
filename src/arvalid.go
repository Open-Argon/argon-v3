package main

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

func ArValidToHash(a any) (any, ArErr) {
	switch a := a.(type) {
	case ArObject:
		if callable, ok := a.obj["__hash__"]; ok {
			value, err := runCall(call{
				Callable: callable,
				Args:     []any{},
			}, stack{}, 0)
			if err.EXISTS {
				return nil, err
			}
			return value, ArErr{}
		}
	}
	return a, ArErr{}
}
