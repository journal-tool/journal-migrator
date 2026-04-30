package lib

func Mapper[T any, O any](slice []T, mapper func(thing T) O) []O {
	var values []O

	for _, object := range slice {
		value := mapper(object)
		values = append(values, value)
	}

	return values
}

func Intersect[T comparable](slices ...[]T) []T {
	if len(slices) == 0 {
		return nil
	}

	var sets []map[T]bool
	for _, slice := range slices {
		set := BuildSet(slice)
		sets = append(sets, set)
	}

	var result []T
	for _, val := range slices[0] {
		common := true

		for _, set := range sets {
			_, ok := set[val]
			if !ok {
				common = false
				break
			}
		}

		if common {
			result = append(result, val)
		}
	}

	return result
}
