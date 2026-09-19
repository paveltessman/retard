package productexceptself

func productExceptSelf(nums []int) []int {
	answer := make([]int, len(nums))
	answer[0] = 1

	leftProduct := 1
	for i := 1; i < len(nums); i++ {
		leftProduct *= nums[i-1]
		answer[i] = leftProduct
	}

	rightProduct := 1
	for i := len(nums) - 1; i >= 0; i-- {
		answer[i] *= rightProduct
		rightProduct *= nums[i]
	}
	return answer
}
