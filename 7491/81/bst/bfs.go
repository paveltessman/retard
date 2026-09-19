package bst

func deq(q *[]*node) *node {
	n := (*q)[0]
	*q = (*q)[1:]
	return n
}

func levelOrder(root *node) [][]int {
	values := make([][]int, 0)

	if root == nil {
		return values
	}

	queue := []*node{root}
	for len(queue) > 0 {
		size := len(queue)
		level := make([]int, 0, size)

		for range size {
			n := deq(&queue)
			level = append(level, n.val)

			if n.left != nil {
				queue = append(queue, n.left)
			}
			if n.right != nil {
				queue = append(queue, n.right)
			}
		}

		values = append(values, level)
	}
	return values
}
