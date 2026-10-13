package productexceptself

func productExceptSelf(nums []int) []int {

	prefixes := make([]int, len(nums))
	suffixes := make([]int, len(nums))
	result := make([]int, len(nums))

	for i := range nums {
		if i == 0 {
			prefixes[i] = 1
			continue
		}
		prefixes[i] = prefixes[i-1] * nums[i-1]
	}

	for i := len(nums) - 1; i >= 0; i-- {
		if i == len(nums) - 1 {
			suffixes[i] = 1
			continue
		}
		suffixes[i] = suffixes[i+1] * nums[i+1]
	}

	result[0] = suffixes[0]
	result[len(nums)-1] = prefixes[len(nums)-1]
	for i := 1; i < len(nums); i++ {
		result[i] = prefixes[i] * suffixes[i]
	}

	return result
}
