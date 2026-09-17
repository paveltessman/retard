package ll

func has_cycle[T any](first *node[T]) bool {
	var fast, slow = first, first

	for fast != nil && fast.next != nil {
		slow = slow.next
		fast = fast.next.next
		if slow == fast {
			return true
		}
	}
	return false
}
