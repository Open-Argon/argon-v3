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
	case Callable:
		return "function"
	case builtinFunc:
		return "function"
	case ArObject:
		return x.TYPE
	case accessVariable:
		return "variable"
	}
	return "unknown"
}
