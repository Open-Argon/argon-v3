package main

func anyToBool(x any) bool {
	switch x := x.(type) {
	case string:
		return x != ""
	case number:
		return x.Cmp(newNumber()) != 0
	case bool:
		return x
	case nil:
		return false
	case ArMap:
		return len(x) != 0
	case builtinFunc:
		return true
	case Callable:
		return true
	case ArClass:
		return true
	default:
		return true
	}
}
