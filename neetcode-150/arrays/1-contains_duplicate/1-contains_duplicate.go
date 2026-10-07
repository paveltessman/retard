package containsduplicate

func containsDuplicate(nums []int) bool {
	freqs := make(map[int]int) // value: count

	for _, num := range nums {
		freqs[num] += 1
		if freqs[num] > 1 {
			return true
		}
	}
	return false
}
