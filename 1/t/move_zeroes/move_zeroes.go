package movezeroes

import "slices"

func moveZeroesNaive(nums []int) {

	for i, val := range slices.Backward(nums) {
		if val != 0 {
			continue
		}
		for ; i < len(nums)-1 && nums[i+1] != 0; i++ {
			nums[i], nums[i+1] = nums[i+1], nums[i]
		}
	}
}

func moveZeroes(nums []int) {
	slow := 0
	for fast, val := range nums {
		if val == 0 {
			continue
		}
		nums[slow], nums[fast] = nums[fast], nums[slow]
		slow++
	}
}
