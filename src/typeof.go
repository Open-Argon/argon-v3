package main

func typeof(val any) string {
	switch val.(type) {
	case number:
		return "number"
	case nil:
		return "null"
	case bool:
		return "boolean"
	case string:
		return "string"
	case anymap:
		return "array"
	case Callable:
		return "function"
	case builtinFunc:
		return "function"
	case ArObject:
		return "map"
	case accessVariable:
		return "variable"
	}
	return "unknown"
}
