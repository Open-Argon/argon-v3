package main

func typeof(val any) string {
	switch x := val.(type) {
	case number:
		return "number"
	case string:
		return "string"
	case nil:
		return "null"
	case bool:
		return "boolean"
	case Callable:
		return "function"
	case builtinFunc:
		return "function"
	case ArObject:
		if x.TYPE == "array" {
			return "array"
		}
		return "map"
	case accessVariable:
		return "variable"
	}
	return "unknown"
}
