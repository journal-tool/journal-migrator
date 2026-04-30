package lib

func BuildSet[T comparable](slice []T) map[T]bool {
	var set = make(map[T]bool)

	for _, val := range slice {
		set[val] = true
	}

	return set
}
