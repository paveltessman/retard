package searchrotated

func searchStart(nums []int, lo, hi int) int {
	if lo >= hi {
		return lo
	}

	mid := lo + (hi-lo)/2
	if nums[mid] > nums[hi] {
		return searchStart(nums, mid+1, hi)
	}
	return searchStart(nums, lo, mid)
}

func startIndex(nums []int) int {
	if len(nums) < 2 {
		return 0
	}
	return searchStart(nums, 0, len(nums)-1)
}

func at(i, start, size int) int {
	return (start + i) % size
}

func search(nums []int, target int) int {

	start := startIndex(nums)
	size := len(nums)

	lo, hi := 0, len(nums)

	for lo < hi {
		mid := lo + (hi-lo)/2

		if nums[at(mid, start, size)] >= target {
			hi = mid
		} else {
			lo = mid + 1
		}
	}

	if lo < len(nums) && nums[at(lo, start, size)] == target {
		return at(lo, start, size)
	}
	return -1
}
