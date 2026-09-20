package largestrect

import (
	"math"
)

func area(heights []int) int {
	minVal := math.MaxInt
	for _, val := range heights {
		minVal = min(minVal, val)
	}
	return minVal * len(heights)
}

func largestRectangleAreaBrute(heights []int) int {
	largest := 0
	for i := range heights {
		for j := i + 1; j <= len(heights); j++ {
			largest = max(largest, area(heights[i:j]))
		}

	}
	return largest
}

func pop(stack *[]int) int {
	val := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]
	return val
}

func largestRectangleArea(heights []int) int {
	stack := []int{}
	largest := 0

	for i := 0; i <= len(heights); i++ {

		current := 0
		if i < len(heights) {
			current = heights[i]
		}

		for len(stack) > 0 {

			left := pop(&stack)

			if heights[left] <= current {
				stack = append(stack, left)
				break
			}

			height := heights[left]
			width := i
			if len(stack) > 0 {
				width = i - stack[len(stack)-1] - 1
			}

			largest = max(largest, height*width)
		}
		stack = append(stack, i)
	}
	return largest
}
