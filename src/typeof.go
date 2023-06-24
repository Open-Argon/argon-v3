package main

func typeof(val any) string {
	switch x := val.(type) {
	case number:
		return "number"
	case nil:
		return "null"
	case bool:
		return "boolean"
	case string:
		return "string"
	case []any:
		return "array"
	case anymap:
		return "map"
	case Callable:
		return "function"
	case builtinFunc:
		return "function"
	case byte:
		return "byte"
	case []byte:
		return "buffer"
	case ArObject:
		if val, ok := x.obj["__name__"]; ok {
			val := ArValidToAny(val)
			if val, ok := val.(string); ok {
				return val
			}
		}
		return "object"
	case accessVariable:
		return "variable"
	}
	return "unknown"
}
