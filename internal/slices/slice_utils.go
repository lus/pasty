package slices

func Contains[T comparable](src []T, val T) bool {
	for _, elem := range src {
		if elem == val {
			return true
		}
	}
	return false
}
