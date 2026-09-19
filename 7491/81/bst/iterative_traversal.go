package bst

func pop(stack *[]*node) *node {
	n := (*stack)[len(*stack)-1]
	*stack = (*stack)[:len(*stack)-1]
	return n
}

func preorderIterative(root *node) []int {
	values := make([]int, 0)

	if root == nil {
		return values
	}

	stack := []*node{root}

	for len(stack) > 0 {
		current := pop(&stack)
		values = append(values, current.val)

		if current.right != nil {
			stack = append(stack, current.right)
		}

		if current.left != nil {
			stack = append(stack, current.left)
		}
	}
	return values
}

func inorderIterative(root *node) []int {
	values := make([]int, 0)

	stack := make([]*node, 0)

	current := root
	for current != nil || len(stack) > 0 {
		for current != nil {
			stack = append(stack, current)
			current = current.left
		}

		current = pop(&stack)

		values = append(values, current.val)
		current = current.right
	}
	return values
}

func postorderIterative(root *node) []int {
	values := make([]int, 0)

	stack := make([]*node, 0)

	var prev *node
	current := root
	for current != nil || len(stack) > 0 {
		for current != nil {
			stack = append(stack, current)
			current = current.left
		}

		top := pop(&stack)

		if top.right != nil && top.right != prev {
			stack = append(stack, top)
			current = top.right
			continue
		}
		values = append(values, top.val)
		prev = top
	}
	return values
}
