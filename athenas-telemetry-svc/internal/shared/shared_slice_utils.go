package shared

// ================================================================================================

// SliceMap runs the given transformFn on each element in the target slice, allocating a new slice
// the transformed elements.
//
// Example:
//
//	// Output: []int{1, 2, 3}
//	shared.SliceMap([]string{"1", "2", "3"}, strToIntFunc)
//
// Note: If either of the args are nil (or if the input slice is empty), nil will be returned
func SliceMap[T, U any](slice []T, transformFn func(T, int) U) []U {
	if len(slice) == 0 || transformFn == nil {
		return nil
	}
	result := make([]U, len(slice))

	for idx, val := range slice {
		result[idx] = transformFn(val, idx)
	}

	return result
}
