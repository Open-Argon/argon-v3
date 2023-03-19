package main

func AnyToArValid(arr any) any {
	switch arr := arr.(type) {
	case []any:
		return ArArray(arr)
	case string:
		return ArString(arr)
	default:
		return arr
	}
}

func ArValidToAny(a any) any {
	switch a := a.(type) {
	case ArObject:
		switch a.TYPE {
		case "string":
			return a.obj["__value__"]
		case "array":
			return a.obj["__value__"]
		case "class":
			return a.obj["__value__"]
		default:
			return a.obj
		}
	default:
		return a
	}
}
