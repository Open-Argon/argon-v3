package main

func AnyToArValid(arr any) any {
	switch arr := arr.(type) {
	case []any:
		return ArArray(arr)
	case string:
		return ArString(arr)
	case anymap:
		return Map(arr)
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
