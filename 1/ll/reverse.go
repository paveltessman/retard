package ll

func reverse[T any](first *node[T]) *node[T] {
	var prev *node[T]
	var next *node[T]
	curr := first

	for curr != nil {
		next = curr.next
		curr.next = prev
		prev = curr
		curr = next
	}
	return prev
}
