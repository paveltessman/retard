package dailytemperatures

func pop(stack *[]int) int {
	val := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]
	return val
}

func dailyTemperatures(temperatures []int) []int {
	result := make([]int, len(temperatures))
	unresolved := []int{}

	for i, val := range temperatures {
		for len(unresolved) > 0 {
			j := pop(&unresolved)

			if temperatures[j] >= val {
				unresolved = append(unresolved, j)
				break
			}
			result[j] = i - j
		}
		unresolved = append(unresolved, i)
	}
	return result
}
