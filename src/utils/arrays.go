package utils

func Contains(elems []string, v interface{}) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func Filter[T comparable](elems []T, filterFunc func(elem T) bool) []T {
	result := []T{}
	for i := range elems {
		if filterFunc(elems[i]) {
			result = append(result, elems[i])
		}
	}

	return result
}


func Find[T comparable](elems []T, filterFunc func(elem T) bool) T {
	for i := range elems {
		if filterFunc(elems[i]) {
			return elems[i]
		}
	}
	return *new(T)
}