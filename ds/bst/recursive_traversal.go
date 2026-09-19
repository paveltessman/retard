package bst

func preorderRecursive(root *node) []int {
	values := make([]int, 0)

	var visit func(*node)
	visit = func(n *node) {
		if n == nil {
			return
		}

		values = append(values, n.val)
		visit(n.left)
		visit(n.right)
	}
	visit(root)
	return values
}

func inorderRecursive(root *node) []int {
	values := make([]int, 0)

	var visit func(*node)
	visit = func(n *node) {
		if n == nil {
			return
		}

		visit(n.left)
		values = append(values, n.val)
		visit(n.right)
	}
	visit(root)
	return values
}

func postorderRecursive(root *node) []int {
	values := make([]int, 0)

	var visit func(*node)
	visit = func(n *node) {
		if n == nil {
			return
		}

		visit(n.left)
		visit(n.right)
		values = append(values, n.val)
	}
	visit(root)
	return values
}
