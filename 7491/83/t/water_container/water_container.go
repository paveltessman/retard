package watercontainer

func area(height []int, left, right int) int {
	return (right - left) * min(height[right], height[left])
}

func maxArea(height []int) int {
	left := 0
	right := len(height) - 1

	best := 0

	for left < right {
		s := area(height, left, right)
		best = max(best, s)

		if height[left] < height[right] {
			left++
		} else {
			right--
		}
	}
	return best
}
