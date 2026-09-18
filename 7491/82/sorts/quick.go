package sorts

import (
	"cmp"
	"slices"
)

func partHoare[T cmp.Ordered](nums []T, lo, hi int) int {
	pivot := nums[lo]
	i, j := lo-1, hi+1

	for {
		for {
			i++
			if nums[i] >= pivot {
				break
			}
		}
		for {
			j--
			if nums[j] <= pivot {
				break
			}
		}
		if i >= j {
			return j
		}
		nums[i], nums[j] = nums[j], nums[i]
	}
}

func qsort[T cmp.Ordered](nums []T, lo, hi int) {
	if lo >= hi {
		return
	}

	p := partHoare(nums, lo, hi)

	qsort(nums, lo, p)
	qsort(nums, p+1, hi)
}

func quickSort[T cmp.Ordered](nums []T) []T {
	result := slices.Clone(nums)
	qsort(result, 0, len(nums)-1)
	return result
}
