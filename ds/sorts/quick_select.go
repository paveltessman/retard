package sorts

func partition(nums []int, lo, hi int) int {
	pivot := nums[hi]
	i := lo

	for j := lo; j < hi; j++ {
		if nums[j] < pivot {
			nums[i], nums[j] = nums[j], nums[i]
			i++
		}
	}
	nums[i], nums[hi] = nums[hi], nums[i]
	return i
}

func qselect(nums []int, lo, hi, k int) int {
	if lo == hi {
		return nums[lo]
	}

	p := partition(nums, lo, hi)

	if p == k {
		return nums[p]
	}
	if p > k {
		return qselect(nums, lo, p-1, k)
	}
	return qselect(nums, p+1, hi, k)
}

func quickSelect(nums []int, k int) int {
	return qselect(nums, 0, len(nums)-1, k)
}
