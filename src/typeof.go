package main

func typeof(val any) string {
	switch val.(type) {
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
	}
	return "unknown"
}
