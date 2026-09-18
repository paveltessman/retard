package containsduplicate


func containsDuplicate(nums []int) bool {
	heatmap := make(map[int]bool)
	for _, val := range nums {
		if _, ok := heatmap[val]; ok {
			return true
		}
		heatmap[val] = true
	}
	return false
}
