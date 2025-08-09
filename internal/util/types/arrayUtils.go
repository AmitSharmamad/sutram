package types

// RemoveDuplicatesInt32 removes duplicate values from a slice of int32 and returns a new slice
func RemoveDuplicatesInt32(arr []int32) []int32 {
	seen := make(map[int32]struct{})
	result := make([]int32, 0, len(arr))
	for _, v := range arr {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}
