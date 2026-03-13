package misc

func UniqueValues(input ...int64) []int64 {
	uniqueMap := make(map[int64]struct{}) // Map to store unique values

	// Iterate through the input slice and add unique elements to the map
	for _, val := range input {
		uniqueMap[val] = struct{}{}
	}

	// Create a slice to hold the unique values
	uniqueValues := make([]int64, 0, len(uniqueMap))

	for key := range uniqueMap {
		uniqueValues = append(uniqueValues, key)
	}

	return uniqueValues
}

func PtrToVal[T any](value *T, defaultValue T) T {
	if value == nil {
		return defaultValue
	}

	return *value
}
