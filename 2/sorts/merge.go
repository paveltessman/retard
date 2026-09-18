package sorts

import (
	"cmp"
	"slices"
)

func merge[T cmp.Ordered](a, b []T) []T {
	result := make([]T, 0, (len(a) + len(b)))
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		if a[i] < b[j] {
			result = append(result, a[i])
			i++
		} else {
			result = append(result, b[j])
			j++
		}
	}
	result = append(result, a[i:]...)
	result = append(result, b[j:]...)
	return result
}

func mergeSort[T cmp.Ordered](nums []T) []T {
	if len(nums) <= 1 {
		return slices.Clone(nums)
	}

	mid := len(nums) / 2
	left := mergeSort(nums[:mid])
	right := mergeSort(nums[mid:])

	return merge(left, right)
}
