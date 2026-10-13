package longestconsecutive

func longestConsecutive(nums []int) int {

	values := make(map[int]struct{})
	best := 0

	for _, val := range nums {
		values[val] = struct{}{}
	}

	for val := range values {
		if _, seen := values[val-1]; seen {
			continue
		}

		length := 1
		for {
			if _, seen := values[val + length]; !seen {
				best = max(best, length)
				break
			}
			length++
		}
	}
	return best
}
