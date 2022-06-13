package utils

import "log"

func Contains(elems []string, v interface{}) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func Filter[T comparable](elems []T, filterFunc func(elem T) bool) []T {
	log.Print("before", elems)

	result := []T{}
	for i := range elems {
		if filterFunc(elems[i]) {
			result = append(result, elems[i])
		}
	}

	log.Print("after", result)

	return result
}
