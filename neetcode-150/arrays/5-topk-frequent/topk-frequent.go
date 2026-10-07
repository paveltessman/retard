package topkfrequent

func topKFrequent(nums []int, k int) []int {
	freq := make(map[int]int)

	for _, num := range nums {
		freq[num]++
	}

	counts := make([][]int, len(nums)+1)

	for num, count := range freq {
		counts[count] = append(counts[count], num)
	}

	result := make([]int, 0, k)

	for i := len(counts) - 1; i >= 0; i-- {
		nums := counts[i]
		for _, num := range nums {
			result = append(result, num)
			if len(result) == k {
				return result
			}
		}
	}
	panic("unreachable")
}
