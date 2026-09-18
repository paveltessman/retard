package minimumsubarray

func minSubArrayLen(target int, nums []int) int {
	minLen := len(nums) + 1

	left := 0
	sum := 0

	for right, val := range nums {
		sum += val

		for sum >= target {
			windowLen := right - left + 1
			if windowLen < minLen {
				minLen = windowLen
			}

			sum -= nums[left]
			left++
		}
	}
	if minLen == len(nums)+1 {
		return 0
	}
	return minLen
}
