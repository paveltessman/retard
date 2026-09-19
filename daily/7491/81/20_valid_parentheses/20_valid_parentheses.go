package validparentheses

var table = map[rune]rune{
	')': '(',
	']': '[',
	'}': '{',
}

func pop(stack *[]rune) rune {
	r := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]
	return r
}

func isValid(s string) bool {

	stack := make([]rune, 0)

	for _, bracket := range s {
		expected, closing := table[bracket]
		if !closing {
			stack = append(stack, bracket)
			continue
		}

		if len(stack) == 0 {
			return false
		}

		pair := pop(&stack)
		if pair != expected {
			return false
		}
	}
	return len(stack) == 0
}
