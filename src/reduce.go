package main

func reduce[T any](reducer func(x T, y T) T, arr []T) T {
	result := arr[0]
	for i := 1; i < len(arr); i++ {
		result = reducer(arr[i], result)
	}
	return result
}
