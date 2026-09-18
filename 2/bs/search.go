package bs

func binarySearch(nums []int, val int) (int, bool) {
	lo, hi := 0, len(nums)-1

	for lo <= hi {
		mid := lo + (hi-lo)/2

		if nums[mid] == val {
			return mid, true
		}
		if nums[mid] < val {
			lo = mid + 1
			continue
		}
		hi = mid - 1
	}
	return 0, false
}
