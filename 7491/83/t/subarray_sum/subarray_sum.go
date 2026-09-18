package subarraysum

func subarraySumNaive(nums []int, k int) int {
	count := 0
	for i := range nums {
		sum := 0
		for j := i; j < len(nums); j++ {
			sum += nums[j]
			if sum == k {
				count++
			}
		}
	}
	return count
}

func subarraySum(nums []int, k int) int {
	seen := make(map[int]int) // map[prefixSum]count
	seen[0] = 1

	runningSum := 0
	count := 0

	for _, val := range nums {
		runningSum += val
		need := runningSum - k

		if freq, ok := seen[need]; ok {
			count += freq
		}
		seen[runningSum]++
	}
	return count
}
